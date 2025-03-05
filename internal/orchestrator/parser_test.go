package orchestrator

import (
	"testing"
	"go.uber.org/zap"
)

func TestParseExpression(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal("Failed to initialize logger", err)
	}
	defer logger.Sync()

	parser := NewParser(logger)

	tests := []struct {
		expression string
		want       float64
		wantErr    bool
	}{
		{"2 + 2 * 2", 6, false},    // 2 + (2 * 2) = 6
		{"10 / 2 - 1", 4, false},   // (10 / 2) - 1 = 4
		{"1 / 0", 0, true},         // Деление на ноль
		{"invalid", 0, true},       // Некорректное выражение
	}

	for _, tt := range tests {
		result, err := parser.ParseExpression(tt.expression)
		if tt.wantErr && err == nil {
			t.Errorf("ParseExpression(%q) expected error, got nil", tt.expression)
			continue
		}
		if !tt.wantErr && err != nil {
			t.Errorf("ParseExpression(%q) unexpected error: %v", tt.expression, err)
			continue
		}
		if !tt.wantErr && result != tt.want {
			t.Errorf("ParseExpression(%q) = %v, want %v", tt.expression, result, tt.want)
		}
	}
}

func TestParseTask(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal("Failed to initialize logger", err)
	}
	defer logger.Sync()

	parser := NewParser(logger)

	task := &models.Task{
		ID:         "test-task-1",
		Expression: "2 + 2 * 2",
	}
	result, err := parser.ParseTask(task)
	if err != nil {
		t.Errorf("ParseTask(%v) unexpected error: %v", task, err)
		return
	}
	if result != 6 {
		t.Errorf("ParseTask(%v) = %v, want 6", task, result)
	}
}