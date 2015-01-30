package shell

import (
	"bytes"
	"io"
	"fmt"
)

type OutputStreamer struct {
	buf *bytes.Buffer
}

func NewOutputStreamer() (*OutputStreamer) {
	os := &OutputStreamer{
		buf: bytes.NewBuffer([]byte("")),
	}
	return os
}

func (o *OutputStreamer) Write(p []byte) (n int, err error) {
	if n, err = o.buf.Write(p); err != nil {
		return
	}
	err = o.OutputLines()
	return
}

func (o *OutputStreamer) OutputLines() (err error) {
	for {
		line, err := o.buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Print(line)
	}
	return nil
}

