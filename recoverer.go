package zapmw

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// OptionRecoverer is an option to change error response process.
type OptionRecoverer func(w http.ResponseWriter, r *http.Request, rvr interface{}) (isWriteLog bool)

// Recoverer is a middleware that recovers from panics.
func Recoverer(lvl zapcore.Level, msg string, opts ...OptionRecoverer) func(next http.Handler) http.Handler {
	if len(opts) == 0 {
		opts = []OptionRecoverer{RecovererDefault}
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {

					isWriteLog := true
					for _, o := range opts {
						if !o(w, r, rvr) {
							isWriteLog = false
						}
					}

					if isWriteLog {
						logger := GetZap(r)
						if ce := logger.Check(lvl, msg); ce != nil {
							ce.Write(zap.Any("error", rvr))
						}
					}
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// RecovererDefault writes server error response.
func RecovererDefault(w http.ResponseWriter, _ *http.Request, _ interface{}) (isWriteLog bool) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return true
}
