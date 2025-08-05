package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oliveiracmorais/cep-clima/handler"
	"github.com/oliveiracmorais/cep-clima/model"
	"github.com/stretchr/testify/assert"
)

func TestGetClimaPorCEP_Success(t *testing.T) {
	mockService := &MockCEPClimaService{
		ReturnResponse: &model.SuccessResponse{
			TempC: 25.0,
			TempF: 77.0,
			TempK: 298.0,
		},
		ReturnError: nil,
	}

	handler := &handler.CEPHandler{Service: mockService}

	req := httptest.NewRequest("GET", "/cep/12345678", nil)
	w := httptest.NewRecorder()

	handler.GetClimaPorCEP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response model.SuccessResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 25.0, response.TempC)
	assert.Equal(t, 77.0, response.TempF)
	assert.Equal(t, 298.0, response.TempK)
}

func TestGetClimaPorCEP_InvalidCEP(t *testing.T) {
	mockService := &MockCEPClimaService{
		ReturnError: &CustomError{Message: "invalid zipcode"},
	}

	handler := &handler.CEPHandler{Service: mockService}

	req := httptest.NewRequest("GET", "/cep/1234567", nil)
	w := httptest.NewRecorder()

	handler.GetClimaPorCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response model.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", response.Message)
}

func TestGetClimaPorCEP_CEPNotFound(t *testing.T) {
	mockService := &MockCEPClimaService{
		ReturnError: &CustomError{Message: "can not find zipcode"},
	}

	handler := &handler.CEPHandler{Service: mockService}

	req := httptest.NewRequest("GET", "/cep/99999999", nil)
	w := httptest.NewRecorder()

	handler.GetClimaPorCEP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response model.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", response.Message)
}

func TestGetClimaPorCEP_MethodNotAllowed(t *testing.T) {
	handler := handler.NewCEPHandler()

	req := httptest.NewRequest("POST", "/cep/12345678", nil)
	w := httptest.NewRecorder()

	handler.GetClimaPorCEP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestGetClimaPorCEP_NotFound(t *testing.T) {
	handler := handler.NewCEPHandler()

	req := httptest.NewRequest("GET", "/invalid", nil)
	w := httptest.NewRecorder()

	handler.GetClimaPorCEP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHealthCheck(t *testing.T) {
	handler := handler.NewCEPHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

type MockCEPClimaService struct {
	ReturnResponse *model.SuccessResponse
	ReturnError    error
}

func (m *MockCEPClimaService) GetClimaPorCEP(cep string) (*model.SuccessResponse, error) {
	return m.ReturnResponse, m.ReturnError
}

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}
