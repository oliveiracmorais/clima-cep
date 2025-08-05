package service

import (
	"errors"
	"testing"

	"github.com/oliveiracmorais/cep-clima/model"
	"github.com/oliveiracmorais/cep-clima/service"
	"github.com/stretchr/testify/assert"
)

type MockViaCEPClient struct {
	ReturnAddress *model.ViaCEPResponse
	ReturnError   error
}

func (m *MockViaCEPClient) GetAddress(cep string) (*model.ViaCEPResponse, error) {
	return m.ReturnAddress, m.ReturnError
}

type MockWeatherAPIClient struct {
	ReturnWeather *model.WeatherAPIResponse
	ReturnError   error
}

func (m *MockWeatherAPIClient) GetWeather(city string) (*model.WeatherAPIResponse, error) {
	return m.ReturnWeather, m.ReturnError
}

func TestIsValidCEP(t *testing.T) {
	assert.True(t, service.IsValidCEP("12345678"))
	assert.False(t, service.IsValidCEP("1234567"))
	assert.False(t, service.IsValidCEP("123456789"))
	assert.False(t, service.IsValidCEP("1234567a"))
	assert.False(t, service.IsValidCEP("12345-678"))
	assert.True(t, service.IsValidCEP("01001000"))
}

func TestGetClimaPorCEP_Success(t *testing.T) {
	mockViaCEP := &MockViaCEPClient{
		ReturnAddress: &model.ViaCEPResponse{
			Localidade: "São Paulo",
			UF:         "SP",
		},
		ReturnError: nil,
	}

	mockWeatherAPI := &MockWeatherAPIClient{
		ReturnWeather: &model.WeatherAPIResponse{
			Current: struct {
				TempC     float64 `json:"temp_c"`
				TempF     float64 `json:"temp_f"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			}{
				TempC: 25.0,
			},
		},
		ReturnError: nil,
	}

	service := &service.CEPClimaService{
		ViaCEPClient:     mockViaCEP,
		WeatherAPIClient: mockWeatherAPI,
	}

	result, err := service.GetClimaPorCEP("12345678")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 25.0, result.TempC)
	assert.Equal(t, 77.0, result.TempF)  // 25*1.8+32 = 77
	assert.Equal(t, 298.0, result.TempK) // 25+273 = 298
}

func TestGetClimaPorCEP_InvalidCEP(t *testing.T) {
	service := service.NewCEPClimaService()

	result, err := service.GetClimaPorCEP("1234567")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid zipcode", err.Error())
}

func TestGetClimaPorCEP_CEPNotFound(t *testing.T) {
	mockViaCEP := &MockViaCEPClient{
		ReturnError: errors.New("CEP não encontrado"),
	}

	mockWeatherAPI := &MockWeatherAPIClient{}

	service := &service.CEPClimaService{
		ViaCEPClient:     mockViaCEP,
		WeatherAPIClient: mockWeatherAPI,
	}

	result, err := service.GetClimaPorCEP("99999999")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "can not find zipcode", err.Error())
}

func TestGetClimaPorCEP_WeatherNotFound(t *testing.T) {
	mockViaCEP := &MockViaCEPClient{
		ReturnAddress: &model.ViaCEPResponse{
			Localidade: "Cidade Inexistente",
			UF:         "XX",
		},
		ReturnError: nil,
	}

	mockWeatherAPI := &MockWeatherAPIClient{
		ReturnError: errors.New("cidade não encontrada"),
	}

	service := &service.CEPClimaService{
		ViaCEPClient:     mockViaCEP,
		WeatherAPIClient: mockWeatherAPI,
	}

	result, err := service.GetClimaPorCEP("12345678")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "can not find weather")
}
