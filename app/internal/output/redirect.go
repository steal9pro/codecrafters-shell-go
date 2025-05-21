package output

import "slices"

var redirectOutputSymbols = []string{">", "1>"}

func ParseRedirectIfPresent(args []string) (bool, string) {
	for idx, arg := range args {
		if slices.Contains(redirectOutputSymbols, arg) {
			return true, args[idx+1]
		}
	}

	return false, ""
}
