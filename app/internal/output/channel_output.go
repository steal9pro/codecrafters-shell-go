package output

import (
	"bufio"
	"io"
)

type ChannelOutput struct {
	channel chan string
	closed  bool
}

func (co *ChannelOutput) Print(message string) {
	if co.closed {
		return
	}
	co.channel <- message
}

func (co *ChannelOutput) PrintError(message string) {
	if co.closed {
		return
	}
	co.channel <- message
}

func (co *ChannelOutput) WriteStream(r io.Reader, isError bool) {
	if co.closed {
		return
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		co.channel <- string(line)
	}
}

func (co *ChannelOutput) Close() {
	if !co.closed {
		close(co.channel)
		co.closed = true
	}
}

func (co *ChannelOutput) GetChannel() <-chan string {
	return co.channel
}

func NewChannelOutput() *ChannelOutput {
	return &ChannelOutput{
		channel: make(chan string, 100), // Buffered channel to prevent blocking
		closed:  false,
	}
}
