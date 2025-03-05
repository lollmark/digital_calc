package orchestrator

import (
	"digitalcalc/internal/models"
	"sync"
)

// Storage управляет хранением выражений и задач.
type Storage struct {
	expressions map[string]*models.Expression
	taskQueue   chan models.Task
	mutex       sync.Mutex
}

// NewStorage создает новое хранилище.
func NewStorage() *Storage {
	return &Storage{
		expressions: make(map[string]*models.Expression),
		taskQueue:   make(chan models.Task, 100),
	}
}

// AddExpression добавляет выражение в хранилище.
func (s *Storage) AddExpression(expr *models.Expression) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.expressions[expr.ID] = expr
}

// AddTask добавляет задачу в очередь.
func (s *Storage) AddTask(task models.Task) {
	s.taskQueue <- task
}

// GetExpression возвращает выражение по ID.
func (s *Storage) GetExpression(id string) *models.Expression {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.expressions[id]
}

// GetExpressions возвращает список всех выражений.
func (s *Storage) GetExpressions() []*models.Expression {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := make([]*models.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		result = append(result, expr)
	}
	return result
}

// GetTask возвращает следующую задачу из очереди (блокирует, если задач нет).
func (s *Storage) GetTask() models.Task {
	return <-s.taskQueue
}