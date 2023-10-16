package httpes

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}
func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
	return
}

func LogMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := slog.Default()
		wraped := wrapResponseWriter(w)
		start := time.Now()
		userAgent := r.Header.Get("User-Agent")
		ip := r.RemoteAddr
		host := r.Host
		next.ServeHTTP(wraped, r)
		defer log.Info(
			"log handler",
			"ip", ip,
			"host", host,
			"status", wraped.status,
			"method", r.Method,
			"path", r.URL.Path,
			"elapsed", fmt.Sprintf("%v", time.Since(start).String()),
			"user-agent", userAgent,
		)
	})
}
