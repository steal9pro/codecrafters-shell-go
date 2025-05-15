package args

import (
	"bufio"
	"os"
	"strings"
)

func ParseArgs() (string, []string) {
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	command, argsStr, _ := strings.Cut(input[:len(input)-1], " ")

	args := make([]string, 0)
	if len(argsStr) == 0 {
		return command, args
	}
	var currentArg strings.Builder
	inQuotes := false

	// Split the input into args, handling quotes properly
	for i := 0; i < len(argsStr); i++ {
		ch := argsStr[i]

		if ch == '\'' {
			inQuotes = !inQuotes
			continue
		}

		if !inQuotes && ch == ' ' {
			// If we have accumulated characters and hit a space outside quotes
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}

		currentArg.WriteByte(ch)
	}

	// Add the last argument if there is one
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return command, args
}
