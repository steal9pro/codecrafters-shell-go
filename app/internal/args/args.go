package args

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

var preservedSymbols = []byte{'"', '$', '\\'}

func ParseArgs() (string, []string) {
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	command, argsStr, _ := strings.Cut(input[:len(input)-1], " ")

	args := make([]string, 0)
	if len(argsStr) == 0 {
		return command, args
	}
	var currentArg strings.Builder
	inQuotes := false
	inDoubleQuotes := false
	preserveBackslash := false

	// Split the input into args, handling quotes properly
	for i := 0; i < len(argsStr); i++ {
		ch := argsStr[i]

		if slices.Contains(preservedSymbols, ch) && preserveBackslash {
			currentArg.WriteByte(argsStr[i])
			preserveBackslash = false
			continue
		}

		if ch == '\\' {
			if inQuotes {
				currentArg.WriteByte(argsStr[i])
				i++
				if i < len(argsStr) {
					currentArg.WriteByte(argsStr[i])
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
			if i < len(argsStr) {
				currentArg.WriteByte(argsStr[i])
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
		args = append(args, currentArg.String())
	}

	return command, args
}
