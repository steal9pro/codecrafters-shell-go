package output

import (
	"fmt"
	"io"
	"os"
)

type StandartOutput struct {
	stdout io.Writer
	stderr io.Writer
}

func (so *StandartOutput) Print(message string) {
	fmt.Fprintln(so.stdout, message)
}

func (so *StandartOutput) PrintError(message string) {
	fmt.Fprintln(so.stderr, message)
}

func (so *StandartOutput) WriteStream(r io.Reader, isError bool) {
	writer := so.stdout
	if isError {
		writer = so.stderr
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, werr := writer.Write(buf[:n])
			if werr != nil {
				fmt.Println("Error writing to output:", werr)
				return
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading from stream:", err)
			return
		}
	}
}

func NewOutput() Output {
	return &StandartOutput{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}
