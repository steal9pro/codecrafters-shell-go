package cmds

import (
	"fmt"
	"io"
	"os"
)

var fileName = "history.txt"

type History struct {
	nextIndex int
	file      *os.File
}

func InitHistory() *History {
	flags := os.O_RDWR | os.O_APPEND | os.O_CREATE

	file, err := os.OpenFile(fileName, flags, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	return &History{
		nextIndex: 1,
		file:      file,
	}
}

func (h *History) Run(args []string) {
	// Seek to the beginning of the file to read from start
	_, err := h.file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking to start of file:", err)
		return
	}

	buf := make([]byte, 32*1024)

	for {
		n, err := h.file.Read(buf)

		if n > 0 {
			fmt.Print(string(buf[:n]))
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}
}

func (h *History) Write(input string) error {
	message := fmt.Sprintf("    %v  %s", h.nextIndex, input)
	_, err := fmt.Fprintln(h.file, message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	h.nextIndex++

	return nil
}

func (h *History) Close() error {
	err := h.file.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
	}

	return nil
}
