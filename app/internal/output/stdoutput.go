package output

import (
	"fmt"
	"io"
	"os"
)

type StandartOutput struct {
	isError bool
	output  io.Writer
}

func (so *StandartOutput) Print(message string) {
	fmt.Fprintln(so.output, message)
}

func (so *StandartOutput) PrintError(message string) {
	fmt.Fprintln(so.output, message)
}

func (so *StandartOutput) WriteStream(r io.Reader) {
	writer := so.output

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		// fmt.Printf("get new %d bytes from cmd to standart output\n", n)
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
			// fmt.Println("Error reading from stream:", err)
			return
		}
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
