package reader

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/internal/autocompletition"
	"golang.org/x/term"
)

const (
	KEY_TAB       = 9
	KEY_ENTER     = 13
	KEY_CTRL_J    = 10
	KEY_BACKSPACE = 127
	KEY_ESC       = 27
	KEY_DEL       = 127
)

type StreamReader struct {
	tabPressed    bool
	buffer        strings.Builder
	cursor        int
	trie          *autocompletition.TrieNode
	originalState *term.State
}

func NewStreamReader(trie *autocompletition.TrieNode) *StreamReader {
	return &StreamReader{
		trie: trie,
	}
}

func (r *StreamReader) enableRawMode() error {
	fd := int(os.Stdin.Fd())

	// Save current terminal state
	state, err := term.GetState(fd)
	if err != nil {
		return fmt.Errorf("failed to get terminal state: %v", err)
	}
	r.originalState = state

	// Set terminal to raw mode
	if _, err := term.MakeRaw(fd); err != nil {
		return fmt.Errorf("failed to set terminal to raw mode: %v", err)
	}

	return nil
}

func (r *StreamReader) disableRawMode() error {
	if r.originalState == nil {
		return nil
	}

	fd := int(os.Stdin.Fd())

	// Restore original terminal state
	if err := term.Restore(fd, r.originalState); err != nil {
		return fmt.Errorf("failed to restore terminal state: %v", err)
	}

	return nil
}

func (r *StreamReader) ReadCommand() (string, []string, error) {
	r.buffer.Reset()
	r.cursor = 0

	if err := r.enableRawMode(); err != nil {
		return "", nil, err
	}
	defer r.disableRawMode()

	for {
		char := make([]byte, 1)
		n, err := os.Stdin.Read(char)
		if err != nil {
			return "", nil, err
		}
		if n == 0 {
			continue
		}

		switch char[0] {
		case KEY_CTRL_J:
			fmt.Print("\r\n")
			return r.parseCommand()
		case KEY_ENTER:
			fmt.Print("\r\n")
			return r.parseCommand()
		case KEY_TAB:
			r.handleTabCompletion()
			r.tabPressed = true
		case KEY_BACKSPACE:
			r.handleBackspace()
		case KEY_ESC:
			r.handleEscapeSequence()
		default:
			if char[0] >= 32 && char[0] < 127 { // Printable ASCII characters
				r.handleRegularChar(char[0])
			}
		}
	}
}

func (r *StreamReader) handleRegularChar(ch byte) {
	// Insert character at cursor position
	current := r.buffer.String()
	if r.cursor >= len(current) {
		r.buffer.WriteByte(ch)
		r.cursor++
		fmt.Printf("%c", ch)
	} else {
		// Insert in middle
		before := current[:r.cursor]
		after := current[r.cursor:]
		r.buffer.Reset()
		r.buffer.WriteString(before)
		r.buffer.WriteByte(ch)
		r.buffer.WriteString(after)
		r.cursor++

		// Redraw line from cursor position
		fmt.Printf("%c%s", ch, after)
		// Move cursor back to correct position
		for i := 0; i < len(after); i++ {
			fmt.Print("\b")
		}
	}
}

func (r *StreamReader) handleBackspace() {
	if r.cursor > 0 {
		current := r.buffer.String()
		before := current[:r.cursor-1]
		after := current[r.cursor:]

		r.buffer.Reset()
		r.buffer.WriteString(before)
		r.buffer.WriteString(after)
		r.cursor--

		// Move cursor back, clear to end of line, rewrite, then position cursor
		fmt.Print("\b")
		fmt.Print("\033[K") // Clear to end of line
		fmt.Print(after)
		// Move cursor back to correct position
		for i := 0; i < len(after); i++ {
			fmt.Print("\b")
		}
	}
}

func (r *StreamReader) handleTabCompletion() {
	current := r.buffer.String()
	words := strings.Fields(current)

	if len(words) == 0 {
		completions, _ := r.trie.GetAllWords("")
		r.showCompletions(completions)
		return
	}

	// Get the last word for completion
	lastWord := words[len(words)-1]

	// For the first word (command), use trie completion
	if len(words) == 1 || (len(words) == 1 && !strings.HasSuffix(current, " ")) {
		completions, longestCommon := r.trie.GetAllWords(lastWord)
		if len(completions) == 1 {
			// Single completion - auto-complete
			completion := completions[0]
			r.replaceLastWord(lastWord, completion)
			r.addSpace()
		} else if len(completions) > 1 {
			// Multiple completions - show options
			if longestCommon != "" {
				r.replaceLastWord(lastWord, longestCommon)
			} else {
				switch r.tabPressed {
				case true:
					r.showCompletions(completions)
				case false:
					r.ringBell()
				}
			}

		} else {
			r.ringBell()
		}
	} else {
		// For arguments, we could implement file/directory completion here
		// For now, just show a message
		fmt.Print("\r\n(file completion not implemented)\r\n")
		r.redrawPrompt()
	}
}

func (r *StreamReader) handleEscapeSequence() {
	// Read the next two characters for arrow keys, etc.
	seq := make([]byte, 2)
	os.Stdin.Read(seq)

	if seq[0] == '[' {
		switch seq[1] {
		case 'A': // Up arrow
			// TODO: Implement history navigation
		case 'B': // Down arrow
			// TODO: Implement history navigation
		case 'C': // Right arrow
			if r.cursor < r.buffer.Len() {
				r.cursor++
				fmt.Print("\033[C")
			}
		case 'D': // Left arrow
			if r.cursor > 0 {
				r.cursor--
				fmt.Print("\033[D")
			}
		}
	}
}

func (r *StreamReader) addSpace() {
	r.buffer.WriteRune(' ')
	r.cursor++
	fmt.Print(" ")
}

func (r *StreamReader) ringBell() {
	fmt.Print("\x07")
}

func (r *StreamReader) replaceLastWord(oldWord, newWord string) {
	// Calculate how many characters to backspace
	backspaceCount := len(oldWord)

	// Backspace to remove old word
	for i := 0; i < backspaceCount; i++ {
		fmt.Print("\b \b")
	}

	// Update buffer
	current := r.buffer.String()
	newContent := current[:len(current)-len(oldWord)] + newWord
	r.buffer.Reset()
	r.buffer.WriteString(newContent)
	r.cursor = r.buffer.Len()

	// Print new word
	fmt.Print(newWord)
}

func (r *StreamReader) showCompletions(completions []string) {
	if len(completions) == 0 {
		return
	}

	fmt.Print("\r\n")
	for _, completion := range completions {
		fmt.Printf("%s  ", completion)
	}
	fmt.Print("\r\n")
	r.redrawPrompt()
}

func (r *StreamReader) redrawPrompt() {
	fmt.Print("$ ")
	current := r.buffer.String()
	fmt.Print(current)
	// Position cursor correctly
	for i := r.buffer.Len(); i > r.cursor; i-- {
		fmt.Print("\b")
	}
}

func (r *StreamReader) parseCommand() (string, []string, error) {
	input := strings.TrimSpace(r.buffer.String())
	if input == "" {
		return "", nil, nil
	}

	// Use similar parsing logic as the original args package
	args := make([]string, 0)
	inQuotes := false
	inDoubleQuotes := false
	preserveBackslash := false

	var currentArg strings.Builder

	for i := 0; i < len(input); i++ {
		ch := input[i]

		preservedSymbols := []byte{'"', '$', '\\'}
		contains := func(slice []byte, item byte) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		}

		if contains(preservedSymbols, ch) && preserveBackslash {
			currentArg.WriteByte(input[i])
			preserveBackslash = false
			continue
		}

		if ch == '\\' {
			if inQuotes {
				currentArg.WriteByte(input[i])
				i++
				if i < len(input) {
					currentArg.WriteByte(input[i])
				}
				continue
			}

			if inDoubleQuotes && !preserveBackslash {
				preserveBackslash = true
				continue
			}
		}

		if ch == '"' && !inQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		if ch == '\'' && !inDoubleQuotes {
			inQuotes = !inQuotes
			continue
		}

		if ch == '\\' && !inQuotes && !inDoubleQuotes {
			i++
			if i < len(input) {
				currentArg.WriteByte(input[i])
			}
			continue
		}

		if !inQuotes && !inDoubleQuotes && ch == ' ' {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}

		if preserveBackslash {
			currentArg.WriteByte('\\')
			preserveBackslash = false
		}

		currentArg.WriteByte(ch)
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	if len(args) == 0 {
		return "", nil, nil
	}

	command := args[0]
	cmdArgs := args[1:]

	return command, cmdArgs, nil
}
