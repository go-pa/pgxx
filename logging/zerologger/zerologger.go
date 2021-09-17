package zerologger

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

// DebugContextLogger gets the zerolog logger from the context and returns a debug level logger
func DebugContextLogger(ctx context.Context) io.Writer {
	logger := zerolog.Ctx(ctx)
	return wrapper{logger.Level(zerolog.DebugLevel)}
}

type wrapper struct{ zerolog.Logger }

// Write implements the io.Writer interface.
func (l wrapper) Write(p []byte) (n int, err error) {
	// n = len(p)
	// if n > 0 && p[n-1] == '\n' {
	// 	// Trim CR added by stdlog.
	// 	p = p[0 : n-1]
	// }
	l.Log().CallerSkipFrame(1).Msg(string(p))
	return
}
