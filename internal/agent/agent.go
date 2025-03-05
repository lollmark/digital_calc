package agent

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"digitalcalc/internal/calculator"
	"digitalcalc/internal/models"
)

func Work(logger *zap.Logger) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Warn("Failed to get task", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		logger.Error("Failed to decode task", zap.Error(err))
		return
	}

	result, err := calculator.Calc(task.Expression)
	if err != nil {
		logger.Error("Failed to evaluate expression", zap.String("expression", task.Expression), zap.Error(err))
		return
	}

	postResult(logger, task.ID, result)
}

func postResult(logger *zap.Logger, id string, result float64) {
	data := struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}{ID: id, Result: result}
	body, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error("Failed to post result", zap.String("id", id), zap.Error(err))
		return
	}
	logger.Info("Result posted successfully", zap.String("id", id), zap.Float64("result", result))
}
