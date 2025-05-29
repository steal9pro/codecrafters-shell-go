package output

import (
	"io"
)

type Output interface {
	Print(message string)
	PrintError(message string)
	WriteStream(r io.Reader)
}
