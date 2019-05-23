package zapmw

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

func TestRequest(t *testing.T) {
	withLogger(t, zapcore.InfoLevel, nil, func(logger *zap.Logger, logs *observer.ObservedLogs) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		mw := newMws(
			WithZap(logger),
			Request(zapcore.InfoLevel, "request"),
		)
		r := http.NewServeMux()

		r.Handle("/", mw.then(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		}))
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal("Response Code should be 200")
		}

		obs := logs.AllUntimed()
		if len(obs) != 1 {
			t.Fatal("Log writing should only be once")
		}
		if obs[0].Level != zapcore.InfoLevel {
			t.Fatal("Log writing level should be INFO")
		}
		if obs[0].Message != "request" {
			t.Fatal(`Log writing message should be "request"`)
		}

		wants := []zap.Field{zap.String("method", "GET"),
			zap.String("url", "/"),
			zap.String("proto", "HTTP/1.1"),
			zap.Int("status", 200),
			zap.String("ip", ""),
			zap.Int("byte", len([]byte("Hello World"))),
			zap.Duration("took", time.Second),
		}
		if obs[0].Context[0] != wants[0] {
			t.Fatalf(`Log writing field[0] = %v, want %v`, obs[0].Context[0], wants[0])
		}
		if obs[0].Context[1] != wants[1] {
			t.Fatalf(`Log writing field[1] = %v, want %v`, obs[0].Context[1], wants[1])
		}
		if obs[0].Context[2] != wants[2] {
			t.Fatalf(`Log writing field[2] = %v, want %v`, obs[0].Context[2], wants[2])
		}
		if obs[0].Context[3] != wants[3] {
			t.Fatalf(`Log writing field[3] = %v, want %v`, obs[0].Context[3], wants[3])
		}
		if obs[0].Context[4] != wants[4] {
			t.Fatalf(`Log writing field[4] = %v, want %v`, obs[0].Context[4], wants[4])
		}
		if obs[0].Context[5] != wants[5] {
			t.Fatalf(`Log writing field[5] = %v, want %v`, obs[0].Context[5], wants[5])
		}
		if obs[0].Context[6].Key != wants[6].Key {
			t.Fatalf(`Log writing field[6].Key = %v, want %v`, obs[0].Context[6].Key, wants[6].Key)
		}
	})
}

func TestRequestWithOption(t *testing.T) {
	withLogger(t, zapcore.InfoLevel, nil, func(logger *zap.Logger, logs *observer.ObservedLogs) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		mw := newMws(
			WithZap(logger),
			Request(zapcore.InfoLevel, "request", func(r *http.Request) []zap.Field {
				return []zap.Field{zap.String("test", "12345")}
			}),
		)
		r := http.NewServeMux()

		r.Handle("/", mw.then(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		}))
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal("Response Code should be 200")
		}

		obs := logs.AllUntimed()
		if len(obs) != 1 {
			t.Fatal("Log writing should only be once")
		}
		if obs[0].Level != zapcore.InfoLevel {
			t.Fatal("Log writing level should be INFO")
		}
		if obs[0].Message != "request" {
			t.Fatal(`Log writing message should be "request"`)
		}

		wants := []zap.Field{zap.String("method", "GET"),
			zap.String("url", "/"),
			zap.String("proto", "HTTP/1.1"),
			zap.Int("status", 200),
			zap.String("ip", ""),
			zap.Int("byte", len([]byte("Hello World"))),
			zap.Duration("took", time.Second),
			zap.String("test", "12345"),
		}
		if obs[0].Context[0] != wants[0] {
			t.Fatalf(`Log writing field[0] = %v, want %v`, obs[0].Context[0], wants[0])
		}
		if obs[0].Context[1] != wants[1] {
			t.Fatalf(`Log writing field[1] = %v, want %v`, obs[0].Context[1], wants[1])
		}
		if obs[0].Context[2] != wants[2] {
			t.Fatalf(`Log writing field[2] = %v, want %v`, obs[0].Context[2], wants[2])
		}
		if obs[0].Context[3] != wants[3] {
			t.Fatalf(`Log writing field[3] = %v, want %v`, obs[0].Context[3], wants[3])
		}
		if obs[0].Context[4] != wants[4] {
			t.Fatalf(`Log writing field[4] = %v, want %v`, obs[0].Context[4], wants[4])
		}
		if obs[0].Context[5] != wants[5] {
			t.Fatalf(`Log writing field[5] = %v, want %v`, obs[0].Context[5], wants[5])
		}
		if obs[0].Context[6].Key != wants[6].Key {
			t.Fatalf(`Log writing field[6].Key = %v, want %v`, obs[0].Context[6].Key, wants[6].Key)
		}
		if obs[0].Context[7].Key != wants[7].Key {
			t.Fatalf(`Log writing field[7] = %v, want %v`, obs[0].Context[7].Key, wants[7].Key)
		}
	})
}
