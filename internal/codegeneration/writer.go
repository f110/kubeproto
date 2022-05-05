package codegeneration

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
)

type Writer struct {
	w *bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{w: new(bytes.Buffer)}
}

func (w *Writer) F(format string, a ...interface{}) {
	if len(format) == 0 {
		_, _ = fmt.Fprintln(w.w)
		return
	}

	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	_, _ = fmt.Fprintf(w.w, format, a...)
}

func (w *Writer) Fn(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w.w, format, a...)
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}

func (w *Writer) WriteTo(out io.Writer) (int64, error) {
	return w.w.WriteTo(out)
}

func (w *Writer) Bytes() []byte {
	return w.w.Bytes()
}

func (w *Writer) String() string {
	return w.w.String()
}

func (w *Writer) Format() error {
	formatted, err := format.Source(w.Bytes())
	if err != nil {
		scanner := bufio.NewScanner(strings.NewReader(w.String()))
		i := 1
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "%d: %s\n", i, scanner.Text())
			i++
		}
		return err
	}
	w.w.Reset()
	w.w.Write(formatted)

	return nil
}

func (w *Writer) DebugOut() {
	scanner := bufio.NewScanner(strings.NewReader(w.w.String()))
	for i := 1; scanner.Scan(); i++ {
		fmt.Fprintf(os.Stderr, "%d: %s\n", i, scanner.Text())
	}
}
