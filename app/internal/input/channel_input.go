package input

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ChannelInput struct {
	channel  <-chan []byte
	buffer   *bytes.Buffer
	scanner  *bufio.Scanner
	closed   bool
	finished bool
}

func NewChannelInput(channel <-chan []byte) io.Reader {
	return &ChannelInput{
		channel:  channel,
		buffer:   &bytes.Buffer{},
		closed:   false,
		finished: false,
	}
}

func (ci *ChannelInput) Read(p []byte) (n int, err error) {
	if ci.finished {
		fmt.Println("finishing")
		return 0, io.EOF
	}

	// If buffer is empty, try to get more data from channel
	if ci.buffer.Len() == 0 {
		select {
		case data, ok := <-ci.channel:
			if !ok {
				ci.finished = true
				return 0, io.EOF
			}
			// Write raw data as-is
			ci.buffer.Write(data)
		default:
			// Channel is empty but not closed
			if ci.closed {
				fmt.Println("ci closed")
				ci.finished = true
				return 0, io.EOF
			}
			// Try one more time with blocking read
			data, ok := <-ci.channel
			if !ok {
				ci.finished = true
				return 0, io.EOF
			}
			ci.buffer.Write(data)
		}
	}

	return ci.buffer.Read(p)
}

func (ci *ChannelInput) Close() {
	ci.closed = true
}

func (ci *ChannelInput) ReadAll() ([]string, error) {
	var chunks []string

	for {
		data, ok := <-ci.channel
		if !ok {
			break
		}
		chunks = append(chunks, string(data))
	}

	return chunks, nil
}

func (ci *ChannelInput) ReadString() (string, error) {
	chunks, err := ci.ReadAll()
	if err != nil {
		return "", err
	}

	return strings.Join(chunks, ""), nil
}

// ConsumeAllChannelInput drains all data from a channel to prevent pipeline hanging
func ConsumeAllChannelInput(channel <-chan []byte) {
	for range channel {
		// Simply consume all data without processing
	}
}
