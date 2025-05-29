package output

import (
	"io"
)

type ChannelOutput struct {
	channel chan []byte
	closed  bool
}

func (co *ChannelOutput) Print(message string) {
	if co.closed {
		return
	}
	co.channel <- []byte(message + "\n")
}

func (co *ChannelOutput) PrintError(message string) {
	if co.closed {
		return
	}
	co.channel <- []byte(message)
}

func (co *ChannelOutput) WriteStream(r io.Reader, isError bool) {
	if co.closed {
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		// fmt.Printf("send new %d bytes from cmd\n", n)
		if n > 0 {
			co.channel <- buf[:n]
		}
		if err != nil {
			break
		}
	}
}

func (co *ChannelOutput) Close() {
	if !co.closed {
		close(co.channel)
		co.closed = true
	}
}

func (co *ChannelOutput) GetChannel() <-chan []byte {
	return co.channel
}

func NewChannelOutput() *ChannelOutput {
	return &ChannelOutput{
		channel: make(chan []byte, 100), // Buffered channel to prevent blocking
		closed:  false,
	}
}
