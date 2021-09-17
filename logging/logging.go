package logging

import (
	"context"
	"io"
	stdlog "log"
)

// LoggerFunc takes an context and returns a io.Writer.A
type LoggerFunc func(ctx context.Context) io.Writer

// DefaultLogger defaults to logging to standard library log.Writer.
// To disable this set it to NoLogger or nil.
var DefaultLogger LoggerFunc = StdLogger

func StdLogger(ctx context.Context) io.Writer {
	return stdlog.Writer()
}

// NoLogger is
func NoLogger(ctx context.Context) io.Writer {
	return nil
}

// Get returns the default logger or nil
func Get(ctx context.Context) io.Writer {
	if DefaultLogger == nil {
		return nil
	}
	return DefaultLogger(ctx)
}
