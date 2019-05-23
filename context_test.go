package zapmw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func withLogger(t testing.TB, e zapcore.LevelEnabler, opts []zap.Option, f func(*zap.Logger, *observer.ObservedLogs)) {
	fac, logs := observer.New(e)
	log := zap.New(fac, opts...)
	f(log, logs)
}

type middleware func(next http.Handler) http.Handler

type mwStack struct {
	middlewares []middleware
}

func newMws(mws ...middleware) mwStack {
	return mwStack{append([]middleware(nil), mws...)}
}

func (m mwStack) then(h func(http.ResponseWriter, *http.Request)) http.Handler {
	var f http.Handler
	f = http.HandlerFunc(h)
	for i := range m.middlewares {
		f = m.middlewares[len(m.middlewares)-1-i](f)
	}
	return f
}

func TestWithZap(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mw := newMws(
		WithZap(logger),
	)
	r := http.NewServeMux()

	var l interface{}
	r.Handle("/", mw.then(func(w http.ResponseWriter, r *http.Request) {
		l = r.Context().Value(ZapKey)
		w.Write([]byte("Hello World"))
	}))
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal("Response Code should be 200")
	}

	if l == nil {
		t.Fatal("Logger in context should'nt be Nil")
	}
	if _, ok := l.(*zap.Logger); !ok {
		t.Fatal("Logger in context should be *zap.Logger instance")
	}
}

func TestGetZap(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mw := newMws(
		WithZap(logger),
	)
	r := http.NewServeMux()

	r.Handle("/", mw.then(func(w http.ResponseWriter, r *http.Request) {
		req = r
		w.Write([]byte("Hello World"))
	}))
	r.ServeHTTP(w, req)

	logger = GetZap(req)
	if logger == nil {
		t.Fatal("Logger in context should'nt be Nil")
	}
}
