package orchestrator

import (
	"digitalcalc/internal/calculator"
	"digitalcalc/internal/models"
	"go.uber.org/zap"
)

// Parser обрабатывает и парсит арифметические выражения.
type Parser struct {
	logger *zap.Logger
}

// NewParser создает новый парсер с логгером.
func NewParser(logger *zap.Logger) *Parser {
	return &Parser{logger: logger}
}

// ParseExpression парсит и вычисляет арифметическое выражение, возвращая результат или ошибку.
func (p *Parser) ParseExpression(expression string) (float64, error) {
	result, err := calculator.Calc(expression)
	if err != nil {
		p.logger.Error("Failed to parse expression", zap.String("expression", expression), zap.Error(err))
		return 0, err
	}
	p.logger.Info("Expression parsed successfully", zap.String("expression", expression), zap.Float64("result", result))
	return result, nil
}

// ParseTask парсит задачу и возвращает результат вычисления.
func (p *Parser) ParseTask(task *models.Task) (float64, error) {
	return p.ParseExpression(task.Expression)
}