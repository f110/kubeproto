package object

import (
	"bytes"
	"fmt"
	"io"
)

type writer struct {
	w *bytes.Buffer
}

func newWriter() *writer {
	return &writer{w: new(bytes.Buffer)}
}

func (w *writer) F(format string, a ...interface{}) {
	if len(format) == 0 {
		_, _ = fmt.Fprintln(w.w)
		return
	}

	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	_, _ = fmt.Fprintf(w.w, format, a...)
}

func (w *writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func (w *writer) WriteTo(out io.Writer) (int64, error) {
	return w.w.WriteTo(out)
}

func (w *writer) Bytes() []byte {
	return w.w.Bytes()
}

func (w *writer) String() string {
	return w.w.String()
}
