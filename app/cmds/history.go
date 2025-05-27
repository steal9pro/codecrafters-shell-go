package cmds

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var fileName = "history.txt"

type History struct {
	nextToWriteIndex int
	file             *os.File
	navigationIndex  int
	lines            []string
}

func InitHistory() *History {
	flags := os.O_RDWR | os.O_CREATE | os.O_TRUNC

	file, err := os.OpenFile(fileName, flags, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	return &History{
		nextToWriteIndex: 1,
		file:             file,
		lines:            make([]string, 0),
	}
}

func (h *History) Run(args []string) {
	if len(args) > 0 {
		var cursor int64

		stat, _ := h.file.Stat()
		filesize := stat.Size()

		rowsAmount, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid number of rows:", err)
			return
		}

		var resp string
		var linesFound int
		for {
			cursor -= 1
			_, err := h.file.Seek(cursor, io.SeekEnd)
			if err != nil {
				fmt.Printf("Error seeking to #%v of file from end: %v", cursor, err)
				return
			}

			char := make([]byte, 1)
			h.file.Read(char)

			if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
				linesFound++
				if linesFound == rowsAmount {
					break
				}
			}
			resp = fmt.Sprintf("%s%s", string(char), resp)

			if cursor == -filesize { // stop if we are at the begining
				break
			}
		}

		fmt.Print(resp)
	} else {
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
}

func (h *History) Write(input string) error {
	message := fmt.Sprintf("    %v  %s", h.nextToWriteIndex, input)
	h.lines = append(h.lines, strings.TrimSpace(input))
	h.navigationIndex++

	_, err := fmt.Fprintln(h.file, message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	h.nextToWriteIndex++

	return nil
}

func (h *History) Close() error {
	err := h.file.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
	}

	return nil
}

func (h *History) Down() string {
	if h.navigationIndex == len(h.lines)-1 {
		return ""
	}

	h.navigationIndex++
	line := h.lines[h.navigationIndex]
	return line
}

func (h *History) Up() string {
	if len(h.lines) == 0 || h.navigationIndex == 0 {
		return ""
	}

	h.navigationIndex--
	line := h.lines[h.navigationIndex]
	return line
}
