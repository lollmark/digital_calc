package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware создает middleware для логирования HTTP-запросов.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаем обёрнутый ResponseWriter для отслеживания статуса
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Записываем время начала обработки запроса
			start := time.Now()

			// Выполняем следующий обработчик
			next.ServeHTTP(ww, r)

			// Логируем запрос
			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// responseWriter обёрнутый http.ResponseWriter для отслеживания статуса ответа.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader переопределяет метод для записи статуса.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}