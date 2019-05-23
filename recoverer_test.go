package zapmw

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverer(t *testing.T) {
	withLogger(t, zapcore.InfoLevel, nil, func(logger *zap.Logger, logs *observer.ObservedLogs) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		mw := newMws(
			WithZap(logger),
			Recoverer(zapcore.ErrorLevel, "recover"),
		)
		r := http.NewServeMux()

		r.Handle("/", mw.then(func(w http.ResponseWriter, r *http.Request) {
			panic("panic")
		}))
		r.ServeHTTP(w, req)

		if w.Code != 500 {
			t.Fatal("Response Code should be 500")
		}

		obs := logs.AllUntimed()
		if len(obs) != 1 {
			t.Fatal("Log writing should only be once")
		}
		if obs[0].Level != zapcore.ErrorLevel {
			t.Fatal("Log writing level should be ERROR")
		}
		if obs[0].Message != "recover" {
			t.Fatal(`Log writing message should be "recover"`)
		}

		wants := []zap.Field{zap.Any("error", "panic")}
		if obs[0].Context[0] != wants[0] {
			t.Fatalf(`Log writing field[0] = %v, want %v`, obs[0].Context[0], wants[0])
		}
	})
}
