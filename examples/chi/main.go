package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/purini-to/zapmw"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var cnt = 1

func withRequestID(logger *zap.Logger, r *http.Request) *zap.Logger {
	reqID := middleware.GetReqID(r.Context())
	if len(reqID) == 0 {
		return logger
	}

	return logger.With(zap.String("reqId", reqID))
}

func main() {
	logger, _ := zap.NewDevelopment()

	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		zapmw.WithZap(logger, withRequestID), // logger with request id.
		zapmw.Request(zapcore.InfoLevel, "request"),
		zapmw.Recoverer(zapcore.ErrorLevel, "recover", zapmw.RecovererDefault),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger := zapmw.GetZap(r)
		logger.Debug("access!", zap.Int("cnt", cnt))
		w.Write([]byte("Hello world"))
		cnt++
	})

	http.ListenAndServe(":3000", r)
}
