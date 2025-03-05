package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"digitalcalc/internal/models"
	"digitalcalc/internal/orchestrator"
)

// CalculateHandler обрабатывает запросы на вычисление арифметического выражения.
func CalculateHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Декодируем тело запроса
		var req struct {
			Expression string `json:"expression"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
			logger.Error("Invalid request payload", zap.Error(err))
			return
		}

		// Генерируем уникальный ID для выражения
		id := orchestrator.GenerateID()
		expression := &models.Expression{
			ID:     id,
			Status: "pending",
		}

		// Добавляем выражение в хранилище оркестратора
		orchestrator.AddExpression(expression)

		// Создаём задачу и добавляем её в очередь
		task := models.Task{
			ID:         id,
			Expression: req.Expression,
		}
		orchestrator.AddTask(task)

		// Формируем и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{ID: id})

		logger.Info("Calculation request received", zap.String("id", id), zap.String("expression", req.Expression))
	}
}