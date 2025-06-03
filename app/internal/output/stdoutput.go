package output

import (
	"fmt"
	"io"
	"os"
)

type StandartOutput struct {
	isError   bool
	output    io.Writer
	lastChar  byte
	hasOutput bool
}

func (so *StandartOutput) Print(message string) {
	fmt.Fprint(so.output, message)
}

func (so *StandartOutput) PrintError(message string) {
	fmt.Fprintln(so.output, message)
}

func (so *StandartOutput) WriteStream(r io.Reader) {
	// Reset state for new command output
	so.lastChar = 0
	so.hasOutput = false

	writer := so.output

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, werr := writer.Write(buf[:n])
			if werr != nil {
				fmt.Println("Error writing to output:", werr)
				return
			}
			// Track the last character written
			so.lastChar = buf[n-1]
			so.hasOutput = true
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
	}

	// Ensure output ends with a newline for proper prompt placement
	if so.hasOutput && so.lastChar != '\n' && !so.isError {
		writer.Write([]byte("\n"))
	}
}

func NewOutput(isError bool) Output {
	if isError {
		return &StandartOutput{
			isError: true,
			output:  os.Stderr,
		}
	}
	return &StandartOutput{
		output: os.Stdout,
	}
}
