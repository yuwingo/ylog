package ylog

import "context"

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services.
type Handler interface {
	Log(context.Context, Level, string, ...interface{})
	Close() error
}

// Handlers a bundle for handler with filter function.
type Handlers struct {
	filters  map[string]struct{}
	handlers []Handler
}

func NewHandlers(handlers ...Handler) *Handlers {
	return &Handlers{handlers: handlers}
}

func (hs *Handlers) Log(ctx context.Context, l Level, format string, args ...interface{}) {
	for _, h := range hs.handlers {
		h.Log(ctx, l, format, args...)
	}
}

func (hs *Handlers) Close() (err error) {
	for _, h := range hs.handlers {
		if err := h.Close(); err != nil {

		}
	}
	return
}
