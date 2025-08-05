package service

import (
	"fmt"
	"strings"

	viacep "github.com/oliveiracmorais/cep-clima/internal/via_cep"
	weatherapi "github.com/oliveiracmorais/cep-clima/internal/weather_api"
	"github.com/oliveiracmorais/cep-clima/model"
)

type CEPClimaServiceInterface interface {
	GetClimaPorCEP(cep string) (*model.SuccessResponse, error)
}

type CEPClimaService struct {
	ViaCEPClient     viacep.Client
	WeatherAPIClient weatherapi.Client
}

func NewCEPClimaService() CEPClimaServiceInterface {
	return &CEPClimaService{
		ViaCEPClient:     viacep.NewViaCEPClient(),
		WeatherAPIClient: weatherapi.NewWeatherAPIClient(),
	}
}

func (s *CEPClimaService) GetClimaPorCEP(cep string) (*model.SuccessResponse, error) {
	fmt.Printf("Serviço: Iniciando consulta para CEP: %s\n", cep)

	if !IsValidCEP(cep) {
		fmt.Printf("Serviço: CEP inválido: %s\n", cep)
		return nil, fmt.Errorf("invalid zipcode")
	}

	address, err := s.ViaCEPClient.GetAddress(cep)
	if err != nil {
		if strings.Contains(err.Error(), "CEP não encontrado") {
			fmt.Printf("Serviço: CEP não encontrado: %s\n", cep)
			return nil, fmt.Errorf("can not find zipcode")
		}
		fmt.Printf("Serviço: Erro ao consultar ViaCEP: %v\n", err)
		return nil, err
	}

	fmt.Printf("Serviço: Endereço encontrado - Cidade: %s, UF: %s\n", address.Localidade, address.UF)

	weather, err := s.WeatherAPIClient.GetWeather(address.Localidade)
	if err != nil {
		fmt.Printf("Serviço: Erro ao obter clima: %v\n", err)
		return nil, fmt.Errorf("can not find weather for city: %v", err)
	}

	tempC := weather.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	fmt.Printf("Serviço: Temperaturas calculadas - C: %.1f, F: %.1f, K: %.1f\n", tempC, tempF, tempK)

	return &model.SuccessResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
		City:  address.Localidade,
	}, nil
}

func IsValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	for _, char := range cep {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
