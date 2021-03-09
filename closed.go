package goqa

import "io"

// Closed is differed closure of Closable items with nil checking
func Closed(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
