package main

import (
	"github.com/purini-to/zapmw"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var cnt = 1

func main() {
	logger, _ := zap.NewDevelopment()

	m := http.NewServeMux()
	m.Handle("/",
		zapmw.WithZap(logger)(zapmw.Request(zapcore.InfoLevel, "request")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := zapmw.GetZap(r)
			logger.Debug("access!", zap.Int("cnt", cnt))
			w.Write([]byte("Hello world"))
			cnt++
		}))),
	)

	s := http.Server{Addr: ":3000", Handler: m}
	s.ListenAndServe()
}
