package args

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

var preservedSymbols = []byte{'"', '$', '\\'}

func ParseArgs() (string, []string) {
	args := make([]string, 0)
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	inQuotes := false
	inDoubleQuotes := false
	preserveBackslash := false

	var currentArg strings.Builder

	// Split the input into args, handling quotes properly
	for i := 0; i < len(input); i++ {
		ch := input[i]

		if slices.Contains(preservedSymbols, ch) && preserveBackslash {
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
			// If we have accumulated characters and hit a space outside quotes
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

	// Add the last argument if there is one
	if currentArg.Len() > 0 {
		trimmedArg := strings.TrimRight(currentArg.String(), "\n")
		args = append(args, trimmedArg)
	}

	command := args[0]
	cmdArgs := args[1:]

	return command, cmdArgs
}
