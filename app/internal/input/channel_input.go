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
		case data, ok := <-ci.channel:
			if !ok {
				ci.finished = true
				return 0, io.EOF
			}
			// Write raw data as-is
			ci.buffer.WriteString(data)
		default:
			// Channel is empty but not closed
			if ci.closed {
				ci.finished = true
				return 0, io.EOF
			}
			// Try one more time with blocking read
			data, ok := <-ci.channel
			if !ok {
				ci.finished = true
				return 0, io.EOF
			}
			ci.buffer.WriteString(data)
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
		chunks = append(chunks, data)
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