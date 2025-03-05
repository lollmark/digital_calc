package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"digitalcalc/internal/models"
	"digitalcalc/internal/orchestrator"
)

func CalculateHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Expression string `json:"expression"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
			logger.Error("Invalid request payload", zap.Error(err))
			return
		}

		id := orchestrator.GenerateID()
		expression := &models.Expression{
			ID:      id,
			RawExpr: req.Expression,
			Status:  "pending",
		}

		storage := orchestrator.NewStorage()
		storage.AddExpression(expression)

		task := models.Task{
			ID:         id,
			Expression: req.Expression,
		}
		storage.AddTask(task) 

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{ID: id})

		logger.Info("Calculation request received", zap.String("id", id), zap.String("expression", req.Expression))
	}
}	
