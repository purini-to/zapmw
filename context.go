package zapmw

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type ctxKeyZap int

// ZapKey is the key that holds the unique zap logger in a request context.
const ZapKey ctxKeyZap = iota

// OptionWithZap is an option to add processing to zap logger held in request context.
type OptionWithZap func(logger *zap.Logger, r *http.Request) *zap.Logger

// WithZap is a middleware that sets the zap logger in a context chain.
func WithZap(logger *zap.Logger, opts ...OptionWithZap) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			l := logger
			for _, o := range opts {
				l = o(l, r)
			}
			r = r.WithContext(context.WithValue(r.Context(), ZapKey, l))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// GetZap returns the zap logger in a context.
func GetZap(r *http.Request) *zap.Logger {
	z, _ := r.Context().Value(ZapKey).(*zap.Logger)
	return z
}
