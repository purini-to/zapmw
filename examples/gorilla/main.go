package main

import (
	"github.com/gorilla/mux"
	"github.com/purini-to/zapmw"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var cnt = 1

func main() {
	logger, _ := zap.NewDevelopment()

	r := mux.NewRouter()

	r.Use(
		zapmw.WithZap(logger),
		zapmw.Request(zapcore.InfoLevel, "request"),
		zapmw.Recoverer(zapcore.ErrorLevel, "recover", zapmw.RecovererDefault),
	)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger := zapmw.GetZap(r)
		logger.Debug("access!", zap.Int("cnt", cnt))
		w.Write([]byte("Hello world"))
		cnt++
	})

	http.ListenAndServe(":3000", r)
}
