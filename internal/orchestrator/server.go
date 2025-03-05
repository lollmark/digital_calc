package orchestrator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"digitalcalc/internal/models"
	"go.uber.org/zap"
)

type Server struct {
	logger      *zap.Logger
	parser      *Parser
	expressions map[string]*models.Expression
	taskQueue   chan models.Task
	mutex       sync.Mutex
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger:      logger,
		parser:      NewParser(logger),
		expressions: make(map[string]*models.Expression),
		taskQueue:   make(chan models.Task, 100),
	}
}

func (s *Server) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		s.logger.Error("Invalid request payload", zap.Error(err))
		return
	}

	id := GenerateID()
	expression := &models.Expression{
		ID:      id,
		RawExpr: req.Expression,
		Status:  "pending",
	}
	s.addExpression(expression)

	task := models.Task{
		ID:         id,
		Expression: req.Expression,
	}
	s.addTask(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{ID: id})
	s.logger.Info("Calculation request received", zap.String("id", id), zap.String("expression", req.Expression))
}

func (s *Server) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	expressions := make([]*models.Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		expressions = append(expressions, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Expressions []*models.Expression `json:"expressions"`
	}{Expressions: expressions})
	s.logger.Info("Expressions list retrieved", zap.Int("count", len(expressions)))
}

func (s *Server) HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if expr, exists := s.expressions[id]; exists {
		json.NewEncoder(w).Encode(map[string]*models.Expression{"expression": expr})
		s.logger.Info("Expression retrieved by ID", zap.String("id", id))
	} else {
		http.Error(w, "Expression not found", http.StatusNotFound)
		s.logger.Warn("Expression not found", zap.String("id", id))
	}
}

func (s *Server) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-s.taskQueue:
		json.NewEncoder(w).Encode(task)
		s.logger.Info("Task dispatched", zap.String("id", task.ID))
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
		s.logger.Warn("No tasks available in queue")
	}
}

func (s *Server) HandlePostTaskResult(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request payload", http.StatusUnprocessableEntity)
		s.logger.Error("Invalid result payload", zap.Error(err))
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if expr, exists := s.expressions[result.ID]; exists {
		expr.Status = models.StatusCompleted // Исправлено с "done" на StatusCompleted
		resultStr := fmt.Sprintf("%f", result.Result)
		expr.Result = &resultStr
		w.WriteHeader(http.StatusOK)
		s.logger.Info("Task result processed", zap.String("id", result.ID), zap.Float64("result", result.Result))
	} else {
		http.Error(w, "Task not found", http.StatusNotFound)
		s.logger.Warn("Task not found", zap.String("id", result.ID))
	}
}

func (s *Server) addExpression(expr *models.Expression) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.expressions[expr.ID] = expr
}

func (s *Server) addTask(task models.Task) {
	s.taskQueue <- task
}

func GenerateID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("exp-%d-%d", time.Now().UnixNano(), rand.Intn(1000000))
}
