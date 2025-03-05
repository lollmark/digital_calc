package router

import (
	"net/http"

	"digitalcalc/internal/middleware"
	"digitalcalc/internal/orchestrator"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) http.Handler {
	server := orchestrator.NewServer(logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", server.HandleCalculate)
	mux.HandleFunc("/api/v1/expressions", server.HandleGetExpressions)
	mux.HandleFunc("/api/v1/expressions/", server.HandleGetExpressionByID)
	mux.HandleFunc("/internal/task", server.HandleGetTask)
	mux.HandleFunc("/internal/task/result", server.HandlePostTaskResult)

	loggedRouter := middleware.LoggingMiddleware(logger)(mux)

	logger.Info("Router initialized with logging middleware")

	return loggedRouter
}
