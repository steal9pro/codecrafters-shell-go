package output

import "slices"

var redirectOutputSymbols = []string{">", "1>"}
var redirectErrSymbols = []string{"2>", "&>"}
var appendOutputSymbols = []string{"1>>", ">>"}

func ParseRedirectIfPresent(args []string) (bool, bool, bool, string) {
	var redirectStdout, redirectStdErr, appendStdout bool
	var fileName string

	for idx, arg := range args {
		if slices.Contains(redirectOutputSymbols, arg) {
			redirectStdout = true
			fileName = args[idx+1]
			return redirectStdout, redirectStdErr, appendStdout, fileName
		}

		if slices.Contains(redirectErrSymbols, arg) {
			redirectStdErr = true
			fileName = args[idx+1]
			return redirectStdout, redirectStdErr, appendStdout, fileName
		}

		if slices.Contains(appendOutputSymbols, arg) {
			appendStdout = true
			fileName = args[idx+1]
		}
	}

	return redirectStdout, redirectStdErr, appendStdout, fileName
}
