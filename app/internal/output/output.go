package output

import (
	"fmt"
	"io"
	"os"
)

type Output interface {
	Print(message string)
	PrintError(message string)
	WriteStream(r io.Reader, isError bool)
}

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

type FileOutput struct {
	fileName string
	file     *os.File
	closed   bool
}

func (fo *FileOutput) Print(message string) {
	if fo.closed {
		fmt.Println("Error: file already closed")
		return
	}
	if fo.file == nil {
		var err error
		fo.file, err = os.Create(fo.fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
	}

	_, err := fmt.Fprintln(fo.file, message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func (fo *FileOutput) PrintError(message string) {
	if fo.closed {
		fmt.Println("Error: file already closed")
		return
	}
	if fo.file == nil {
		var err error
		fo.file, err = os.Create(fo.fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
	}

	_, err := fmt.Fprintln(fo.file, message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func (fo *FileOutput) WriteStream(r io.Reader, isError bool) {
	if fo.closed {
		fmt.Println("Error: file already closed")
		return
	}
	if fo.file == nil {
		var err error
		fo.file, err = os.Create(fo.fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, werr := fo.file.Write(buf[:n])
			if werr != nil {
				fmt.Println("Error writing to file:", werr)
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

func NewFileOutput(filename string) Output {
	return &FileOutput{
		fileName: filename,
		closed:   false,
	}
}
