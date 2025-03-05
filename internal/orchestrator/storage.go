package orchestrator

import (
	"digitalcalc/internal/models"
	"sync"
)

type Storage struct {
	expressions map[string]*models.Expression
	taskQueue   chan models.Task
	mutex       sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		expressions: make(map[string]*models.Expression),
		taskQueue:   make(chan models.Task, 100),
	}
}

func (s *Storage) AddExpression(expr *models.Expression) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.expressions[expr.ID] = expr
}

func (s *Storage) AddTask(task models.Task) {
	s.taskQueue <- task
}

func (s *Storage) GetExpression(id string) *models.Expression {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.expressions[id]
}

func (s *Storage) GetExpressions() []*models.Expression {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := make([]*models.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		result = append(result, expr)
	}
	return result
}

func (s *Storage) GetTask() models.Task {
	return <-s.taskQueue
}
