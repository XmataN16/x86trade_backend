package middleware

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// respWriter обёртка для ResponseWriter — захватываем статус и тело ответа.
type respWriter struct {
	http.ResponseWriter
	status int
	buf    *bytes.Buffer
}

func (rw *respWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *respWriter) Write(b []byte) (int, error) {
	// копируем в буфер (ограничиваем размер логируемого тела до, например, 8KB)
	if rw.buf.Len() < 8*1024 {
		remaining := 8*1024 - rw.buf.Len()
		if len(b) > remaining {
			rw.buf.Write(b[:remaining])
		} else {
			rw.buf.Write(b)
		}
	}
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware возвращает middleware, который логгирует каждый запрос/ответ,
// если debug == true. В лог уходит: method, path, status, duration и первые 8KB тела ответа.
func LoggingMiddleware(debug bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !debug {
			// Возвращаем прозрачный middleware (ничего не делаем)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &respWriter{
				ResponseWriter: w,
				status:         200,
				buf:            &bytes.Buffer{},
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			// метод путь query
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = path + "?" + r.URL.RawQuery
			}

			// Короче, логируем результат. Ограничиваем тело, переводим новые строки
			body := strings.ReplaceAll(rw.buf.String(), "\n", "\\n")
			if len(body) > 1024 {
				body = body[:1024] + "...(truncated)"
			}

			log.Printf("[DEBUG] %s %s -> %d (%s) resp_body=[%s]\n", r.Method, path, rw.status, duration, body)
		})
	}
}

// GetDebugFromEnv удобно читать DEBUG из окружения. Возвращает true если
// переменная равна "1" или "true" (регистронезависимо).
func GetDebugFromEnv() bool {
	v := os.Getenv("DEBUG")
	if v == "" {
		return false
	}
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true"
}
