package router

import (
	"net/http"

	"digitalcalc/internal/middleware"
	"digitalcalc/internal/orchestrator"
	"go.uber.org/zap"
)

// NewRouter создает новый маршрутизатор с применением middleware и обработчиков.
func NewRouter(logger *zap.Logger) http.Handler {
	// Создаем сервер оркестратора
	server := orchestrator.NewServer(logger)

	// Создаем мультиплексор маршрутов
	mux := http.NewServeMux()

	// Регистрируем обработчики эндпоинтов
	mux.HandleFunc("/api/v1/calculate", server.HandleCalculate)
	mux.HandleFunc("/api/v1/expressions", server.HandleGetExpressions)
	mux.HandleFunc("/api/v1/expressions/", server.HandleGetExpressionByID)
	mux.HandleFunc("/internal/task", server.HandleGetTask)
	mux.HandleFunc("/internal/task/result", server.HandlePostTaskResult)

	// Применяем middleware для логирования
	loggedRouter := middleware.LoggingMiddleware(logger)(mux)

	// Логируем инициализацию маршрутизатора
	logger.Info("Router initialized with logging middleware")

	return loggedRouter
}