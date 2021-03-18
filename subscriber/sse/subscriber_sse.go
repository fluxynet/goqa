package sse

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

func init() {
	var _ goqa.Subscriber = New(nil, nil)
}

func New(writer io.Writer, flusher http.Flusher) *SSE {
	return &SSE{writer: writer, flusher: flusher}
}

type SSE struct {
	subscriber.Identifiable
	ctx     context.Context
	flusher http.Flusher
	writer  io.Writer
}

func (s *SSE) Notify(event goqa.Event) error {
	var (
		ev   = strings.ReplaceAll(event.Name(), "\n", "_")
		data = strings.ReplaceAll(event.String(), "\n", "_")
	)

	if ev != "" {
		fmt.Fprintf(s.writer, "event: %s\n", ev)
	}
	fmt.Fprintf(s.writer, "data: %s\n\n", strings.ReplaceAll(data, "\n", "\ndata: "))

	if s.flusher != nil {
		s.flusher.Flush()
	}

	return nil
}

func (s *SSE) Serialize() (string, error) {
	return "", subscriber.ErrSerializeNotSupported
}

func (s *SSE) Unserialize(string) error {
	return subscriber.ErrSerializeNotSupported
}
