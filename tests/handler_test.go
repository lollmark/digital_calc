package tests

import (
	"bytes"
	"encoding/json"
	"testing"
	
	"digitalcalc/internal/handler"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

func TestCalculateHandler(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Не удалось инициализировать логгер: %v", err)
	}
	defer logger.Sync()

	handlerFunc := handler.CalculateHandler(logger)

	tests := []struct {
		name           string
		method         string
		payload        interface{}
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Valid Expression",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "1 + 2 * 3"},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"result": "7"},
		},
		{
			name:           "Invalid Characters",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "1 + a"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrInvalidInput},
		},
		{
			name:           "Division by Zero",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "10 / 0"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrDivisionByZero},
		},
		{
			name:           "Missing Expression Field",
			method:         http.MethodPost,
			payload:        map[string]string{"expr": "1 + 2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrMissingField},
		},
		{
			name:           "Empty Expression",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": ""},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrMissingField},
		},
		{
			name:           "Unsupported HTTP Method",
			method:         http.MethodGet, 
			payload:        nil,            
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   map[string]string{"error": errors.ErrUnsupportedMethod},
		},
		{
			name:           "Malformed JSON",
			method:         http.MethodPost,
			payload:        `{"expression": "1 + 2",`, 
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": errors.ErrMalformedJSON},
		},
		{
			name:           "Expression Too Long",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": generateLongExpression(1001)}, 
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": errors.ErrTooLongExpression},
		},
		{
			name:           "Trigger Internal Server Error by Header",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": "1+1"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Internal Server Error"},
		},
		{
			name:           "Trigger Internal Server Error by Long Expression",
			method:         http.MethodPost,
			payload:        map[string]string{"expression": generateLongExpression(800)},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Expression length triggered server error"},
		},
	}

	for _, tt := range tests {
		tt := tt 
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			switch payload := tt.payload.(type) {
			case string:
				reqBody = []byte(payload)
			case map[string]string:
				reqBody, err = json.Marshal(payload)
				if err != nil {
					t.Fatalf("Не удалось сериализовать payload: %v", err)
				}
			case nil:
				reqBody = nil
			default:
				t.Fatalf("Неподдерживаемый тип payload: %T", payload)
			}

			req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.name == "Trigger Internal Server Error by Header" {
				req.Header.Set("X-Trigger-500", "true")
			}

			rr := httptest.NewRecorder()

			handlerFunc.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, rr.Code)
			}

			var responseBody map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Не удалось декодировать тело ответа: %v", err)
			}

			for key, expectedValue := range tt.expectedBody {
				if value, exists := responseBody[key]; !exists || value != expectedValue {
					t.Errorf("Для ключа '%s' ожидалось '%s', получено '%s'", key, expectedValue, value)
				}
			}
		})
	}
}

func generateLongExpression(length int) string {
	expression := ""
	for i := 0; i < length; i++ {
		expression += "1"
	}
	return expression
}

