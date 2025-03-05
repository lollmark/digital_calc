package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"digitalcalc/internal/middleware"
	"go.uber.org/zap"
)

func TestLoggingMiddleware(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer logger.Sync()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) 
		w.WriteHeader(http.StatusOK)
	})


	loggedHandler := middleware.LoggingMiddleware(logger)(testHandler)


	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()


	loggedHandler.ServeHTTP(w, req)


	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Result().StatusCode)
	}
}
