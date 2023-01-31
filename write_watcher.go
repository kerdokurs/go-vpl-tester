package main

import "bytes"

type writeWatcher struct {
	buf *bytes.Buffer
	out chan string
}

func (tw *writeWatcher) Write(b []byte) (int, error) {
	n, err := tw.buf.Write(b)

	if tw.out != nil {
		tw.out <- string(b)
	}

	return n, err
}
