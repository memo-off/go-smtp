package smtp

import (
	"bufio"
	"fmt"
	"io"
)

// A Writer implements convenience methods for writing
// requests or responses to a text protocol network connection.
type Writer struct {
	W   *bufio.Writer
	dot *dotWriter
}

// NewWriter returns a new Writer writing to w.
func NewWriter(w *bufio.Writer) *Writer {
	return &Writer{W: w}
}

var crnl = []byte{'\r', '\n'}
var dotcrnl = []byte{'.', '\r', '\n'}

// PrintfLine writes the formatted output followed by \r\n.
func (w *Writer) PrintfLine(format string, args ...interface{}) error {
	w.closeDot()
	fmt.Fprintf(w.W, format, args...)
	w.W.Write(crnl)
	return w.W.Flush()
}

// DotWriter returns a writer that can be used to write a dot-encoding to w.
// It takes care of inserting leading dots when necessary,
// translating line-ending \n into \r\n, and adding the final .\r\n line
// when the DotWriter is closed. The caller should close the
// DotWriter before the next call to a method on w.
//
// See the documentation for Reader's DotReader method for details about dot-encoding.
func (w *Writer) DotWriter() io.WriteCloser {
	w.closeDot()
	w.dot = &dotWriter{w: w}
	return w.dot
}

func (w *Writer) closeDot() {
	if w.dot != nil {
		w.dot.Close() // sets w.dot = nil
	}
}

type dotWriter struct {
	w *Writer
}

func (d *dotWriter) Write(b []byte) (n int, err error) {
	bw := d.w.W
	return bw.Write(b)
}

func (d *dotWriter) Close() error {
	if d.w.dot == d {
		d.w.dot = nil
	}
	bw := d.w.W
	bw.Write(crnl)
	bw.Write(dotcrnl)
	return bw.Flush()
}
