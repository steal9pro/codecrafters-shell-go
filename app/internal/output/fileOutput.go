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
	fo := &FileOutput{
		fileName: filename,
		closed:   false,
	}

	file, err := os.Create(fo.fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	fo.file = file

	return fo
}
