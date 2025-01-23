package responder

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MockLogger struct {
	LoggedError string
}

func (m *MockLogger) Error(msg string, fields ...zapcore.Field) {
	m.LoggedError = msg
}

func (m *MockLogger) Info(msg string, fields ...zapcore.Field) {}
func TestNewResponder(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	responder := NewResponder(logger)

	assert.NotNil(t, responder, "Responder должен быть не nil")

	respond, ok := responder.(*Respond)
	assert.True(t, ok, "Тип возвращаемого объекта должен быть *Respond")

	assert.Equal(t, logger, respond.log, "Переданный логгер должен быть установлен в поле log")
}

func TestOutputJSON_Success(t *testing.T) {
	mockLogger := &MockLogger{}
	responder := Respond{log: mockLogger}

	data := map[string]string{"key": "value"}
	expectedJSON, _ := json.Marshal(data)

	recorder := httptest.NewRecorder()
	responder.OutputJSON(recorder, data)

	assert.Equal(t, "application/json;charset=utf-8", recorder.Header().Get("Content-Type"))

	assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	assert.Empty(t, mockLogger.LoggedError, "Не должно быть ошибок логирования")
}

func TestOutputJSON_Error(t *testing.T) {
	mockLogger := &MockLogger{}
	responder := Respond{log: mockLogger}

	invalidData := make(chan int)

	recorder := httptest.NewRecorder()
	responder.OutputJSON(recorder, invalidData)

	assert.Equal(t, "responder json encode error", mockLogger.LoggedError)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestErrorBadRequest_Success(t *testing.T) {
	mockLogger := &MockLogger{}
	responder := Respond{log: mockLogger}

	testError := errors.New("invalid request")

	recorder := httptest.NewRecorder()
	responder.ErrorBadRequest(recorder, testError)

	assert.Equal(t, "application/json;charset=utf-8", recorder.Header().Get("Content-Type"))

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	expectedResponse := Response{
		Success: false,
		Message: testError.Error(),
		Data:    nil,
	}
	var actualResponse Response
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err, "Должно быть корректное JSON-сообщение")
	assert.Equal(t, expectedResponse, actualResponse)

	mockLogger.Error("http response bad request status code")

	assert.Equal(t, "http response bad request status code", mockLogger.LoggedError)
}

func TestErrorInternal_Success(t *testing.T) {
	mockLogger := &MockLogger{}
	responder := Respond{log: mockLogger}

	testError := errors.New("internal server error")

	recorder := httptest.NewRecorder()
	responder.ErrorInternal(recorder, testError)

	assert.Equal(t, "application/json;charset=utf-8", recorder.Header().Get("Content-Type"))

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	expectedResponse := Response{
		Success: false,
		Message: testError.Error(),
		Data:    nil,
	}
	var actualResponse Response
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err, "Должно быть корректное JSON-сообщение")
	assert.Equal(t, expectedResponse, actualResponse)

	assert.Equal(t, "http response internal error", mockLogger.LoggedError)
}
