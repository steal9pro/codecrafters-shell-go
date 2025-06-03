package output

import (
	"io"
	"sync"
)

type PipeOutput struct {
	Writer io.Writer
	mu     sync.Mutex
}

func (po *PipeOutput) Print(message string) {
	po.mu.Lock()
	defer po.mu.Unlock()
	po.Writer.Write([]byte(message))
}

func (po *PipeOutput) PrintError(message string) {
	po.mu.Lock()
	defer po.mu.Unlock()
	po.Writer.Write([]byte(message))
}

func (po *PipeOutput) WriteStream(r io.Reader) {
	po.mu.Lock()
	defer po.mu.Unlock()
	io.Copy(po.Writer, r)
}
