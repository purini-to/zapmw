package zapmw

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// OptionRequest is an option to add fields for request log.
type OptionRequest func(r *http.Request) []zap.Field

// Request a middleware that logs the start and end of each request.
func Request(lvl zapcore.Level, msg string, opts ...OptionRequest) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := NewWrapResponseWriter(w, r.ProtoMajor)
			defer func(begin time.Time) {
				logger := GetZap(r)
				fields := []zap.Field{
					zap.String("method", r.Method),
					zap.String("url", r.URL.String()),
					zap.String("proto", r.Proto),
					zap.Int("status", ww.Status()),
					zap.String("ip", r.RemoteAddr),
					zap.Int("byte", ww.BytesWritten()),
					zap.Duration("took", time.Since(begin)),
				}
				for _, o := range opts {
					fields = append(fields, o(r)...)
				}
				if ce := logger.Check(lvl, msg); ce != nil {
					ce.Write(fields...)
				}
			}(time.Now())
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
