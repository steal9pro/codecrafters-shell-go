package output

import "slices"

var redirectOutputSymbols = []string{">", "1>"}
var redirectErrSymbols = []string{"2>", "&>"}

func ParseRedirectIfPresent(args []string) (bool, bool, string) {
	var redirectStdout, redirectStdErr bool
	var fileName string

	for idx, arg := range args {
		if slices.Contains(redirectOutputSymbols, arg) {
			redirectStdout = true
			fileName = args[idx+1]
			return redirectStdout, redirectStdErr, fileName
		}

		if slices.Contains(redirectErrSymbols, arg) {
			redirectStdErr = true
			fileName = args[idx+1]
			return redirectStdout, redirectStdErr, fileName
		}
	}

	return redirectStdout, redirectStdErr, fileName
}
