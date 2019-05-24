# zapmw

[![CircleCI](https://circleci.com/gh/purini-to/zapmw/tree/master.svg?style=svg)](https://circleci.com/gh/purini-to/zapmw/tree/master)
[![codecov](https://codecov.io/gh/purini-to/zapmw/branch/master/graph/badge.svg)](https://codecov.io/gh/purini-to/zapmw)

zapmw is `net/http.Handler` middleware using [zap](https://github.com/uber-go/zap).  

## Installation

`go get -u github.com/purini-to/zapmw`

## Quick Start

It can be used as middleware on compatible routers due to the net/http.Handler interface.

router http.NewServeMux:
```go
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
// out:
// 2019-05-23T14:00:15.410+0900	DEBUG	cmd/main.go:19	access!	{"cnt": 1}
// 2019-05-23T14:00:15.410+0900	INFO	zapmw/request.go:32	request	{"method": "GET", "url": "/", "proto": "HTTP/1.1", "status": 200, "ip": "[::1]:54434", "byte": 11, "took": "440.597µs"}
// 2019-05-23T14:00:15.917+0900	DEBUG	cmd/main.go:19	access!	{"cnt": 2}
// 2019-05-23T14:00:15.917+0900	INFO	zapmw/request.go:32	request	{"method": "GET", "url": "/", "proto": "HTTP/1.1", "status": 200, "ip": "[::1]:54434", "byte": 11, "took": "86.27µs"}
```

router [chi](https://github.com/go-chi/chi):
```go
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
// out:
// 2019-05-23T14:11:38.650+0900	DEBUG	chi/main.go:37	access!	{"reqId": "xxxxxx-000001", "cnt": 1}
// 2019-05-23T14:11:38.650+0900	INFO	zapmw/request.go:32	request	{"reqId": "xxxxxx-000001", "method": "GET", "url": "/", "proto": "HTTP/1.1", "status": 200, "ip": "[::1]:54517", "byte": 11, "took": "174.504µs"}
// 2019-05-23T14:11:39.442+0900	DEBUG	chi/main.go:37	access!	{"reqId": "xxxxxx-000002", "cnt": 2}
// 2019-05-23T14:11:39.442+0900	INFO	zapmw/request.go:32	request	{"reqId": "xxxxxx-000002", "method": "GET", "url": "/", "proto": "HTTP/1.1", "status": 200, "ip": "[::1]:54517", "byte": 11, "took": "55.185µs"}
```

router [gorilla/mux](https://github.com/gorilla/mux):
```go
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
// out:
// 2019-05-23T14:18:35.966+0900	DEBUG	gorilla/main.go:26	access!	{"cnt": 1}
// 2019-05-23T14:18:35.966+0900	INFO	zapmw/request.go:32	request	{"method": "GET", "url": "/", "proto": "HTTP/1.1", "status": 200, "ip": "[::1]:54563", "byte": 11, "took": "426.231µs"}
// 2019-05-23T14:18:36.588+0900	DEBUG	gorilla/main.go:26	access!	{"cnt": 2}
// 2019-05-23T14:18:36.588+0900	INFO	zapmw/request.go:32	request	{"method": "GET", "url": "/", "proto": "HTTP/1.1", "status": 200, "ip": "[::1]:54563", "byte": 11, "took": "84.474µs"}
```