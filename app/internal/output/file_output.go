package output

import (
	"fmt"
	"io"
	"os"
)

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

	_, err := fmt.Fprintln(fo.file, message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func (fo *FileOutput) WriteStream(r io.Reader) {
	if fo.closed {
		fmt.Println("Error: file already closed")
		return
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

func NewFileOutput(filename string, append bool) Output {
	var file *os.File
	fo := &FileOutput{
		fileName: filename,
		closed:   false,
	}

	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	file, err := os.OpenFile(fo.fileName, flags, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	fo.file = file

	return fo
}
