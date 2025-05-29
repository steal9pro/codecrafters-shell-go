package input

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type ChannelInput struct {
	channel  <-chan string
	buffer   *bytes.Buffer
	scanner  *bufio.Scanner
	closed   bool
	finished bool
}

func NewChannelInput(channel <-chan string) *ChannelInput {
	return &ChannelInput{
		channel: channel,
		buffer:  &bytes.Buffer{},
		closed:  false,
		finished: false,
	}
}

func (ci *ChannelInput) Read(p []byte) (n int, err error) {
	if ci.finished {
		return 0, io.EOF
	}

	// If buffer is empty, try to get more data from channel
	if ci.buffer.Len() == 0 {
		select {
		case line, ok := <-ci.channel:
			if !ok {
				ci.finished = true
				return 0, io.EOF
			}
			// Add newline to match typical pipe behavior
			ci.buffer.WriteString(line + "\n")
		default:
			// Channel is empty but not closed
			if ci.closed {
				ci.finished = true
				return 0, io.EOF
			}
			// Try one more time with blocking read
			line, ok := <-ci.channel
			if !ok {
				ci.finished = true
				return 0, io.EOF
			}
			ci.buffer.WriteString(line + "\n")
		}
	}

	return ci.buffer.Read(p)
}

func (ci *ChannelInput) Close() {
	ci.closed = true
}

func (ci *ChannelInput) ReadAll() ([]string, error) {
	var lines []string
	
	for {
		line, ok := <-ci.channel
		if !ok {
			break
		}
		lines = append(lines, line)
	}
	
	return lines, nil
}

func (ci *ChannelInput) ReadString() (string, error) {
	lines, err := ci.ReadAll()
	if err != nil {
		return "", err
	}
	
	return strings.Join(lines, "\n"), nil
}