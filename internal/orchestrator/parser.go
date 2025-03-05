package orchestrator

import (
	"digitalcalc/internal/calculator"
	"digitalcalc/internal/models"
	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{logger: logger}
}

func (p *Parser) ParseExpression(expression string) (float64, error) {
	result, err := calculator.Calc(expression)
	if err != nil {
		p.logger.Error("Failed to parse expression", zap.String("expression", expression), zap.Error(err))
		return 0, err
	}
	p.logger.Info("Expression parsed successfully", zap.String("expression", expression), zap.Float64("result", result))
	return result, nil
}

func (p *Parser) ParseTask(task *models.Task) (float64, error) {
	return p.ParseExpression(task.Expression)
}
