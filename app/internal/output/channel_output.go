package output

import (
	"bufio"
	"fmt"
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
		fmt.Println("put test to channel")
		line := scanner.Text()
		co.channel <- line
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
