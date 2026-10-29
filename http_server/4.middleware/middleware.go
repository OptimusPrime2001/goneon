package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"time"
)

type Middleware func(http.Handler) http.Handler

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range slices.Backward(mws) { // bọc từ trong ra ngoài
		h = mw(h)
	}
	return h
}

// Loggin
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	n, err := s.ResponseWriter.Write(b)
	s.bytes += n
	return n, err
}

// Cho phép http.NewResponseController tìm ra Flusher/Hijacker bên dưới
func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }

func logging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK} // mặc định 200

			next.ServeHTTP(rec, r)

			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"http", r.ProtoMajor,
				"status", rec.status,
				"bytes", rec.bytes,
				"duration", time.Since(start),
				"request_id", RequestIDFrom(r.Context()),
				"remote", r.RemoteAddr,
			)
		})
	}
}

// Recover
func recoverer(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler {
						panic(rec) // cố ý hủy response, để net/http xử lý
					}
					logger.Error("panic recovered",
						"error", rec,
						"stack", string(debug.Stack()),
						"request_id", RequestIDFrom(r.Context()),
					)
					writeJSON(w, http.StatusInternalServerError,
						map[string]string{"error": "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Request ID:
type ctxKey string // kiểu key riêng, tránh đụng key của package khác

const requestIDKey ctxKey = "request_id"

func newID() string {
	b := make([]byte, 8)
	rand.Read(b) // crypto/rand
	return hex.EncodeToString(b)
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" || len(id) > 64 { // nhận ID từ upstream nếu hợp lý, nếu không thì tự sinh
			id = newID()
		}
		w.Header().Set("X-Request-ID", id) // trả lại cho client

		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx)) // truyền request mới có context mới
	})
}

func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// CORS
func cors(allowedOrigins []string) Middleware {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowed[origin] { // chỉ origin nằm trong whitelist
				h := w.Header()
				h.Set("Access-Control-Allow-Origin", origin)
				h.Add("Vary", "Origin") // cache không được dùng chung kết quả giữa các origin

				// Preflight: trả lời luôn, không chuyển tiếp xuống mux
				if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
					h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
					h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
					h.Set("Access-Control-Max-Age", "600") // trình duyệt nhớ kết quả 10 phút
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
