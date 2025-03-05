package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"

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

	if err := postResult(logger, task.ID, result); err != nil {
		logger.Error("Failed to post result", zap.String("id", task.ID), zap.Error(err))
		return
	}
}

func postResult(logger *zap.Logger, id string, result float64) error {
	data := struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}{ID: id, Result: result}
	body, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal result", zap.String("id", id), zap.Error(err))
		return err
	}
	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error("Failed to post result", zap.String("id", id), zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("Unexpected status code", zap.String("id", id), zap.Int("status", resp.StatusCode))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	logger.Info("Result posted successfully", zap.String("id", id), zap.Float64("result", result))
	return nil
}
