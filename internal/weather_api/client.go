package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/oliveiracmorais/cep-clima/model"
)

const baseURL = "https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no"

type Client interface {
	GetWeather(city string) (*model.WeatherAPIResponse, error)
}

type WeatherAPIClient struct {
	apiKey string
}

func NewWeatherAPIClient() Client {
	apiKey := os.Getenv("WEATHER_API_KEY")
	return &WeatherAPIClient{apiKey: apiKey}
}

func (c *WeatherAPIClient) GetWeather(city string) (*model.WeatherAPIResponse, error) {
	location := fmt.Sprintf("%s,Brazil", city)
	encodedLocation := url.QueryEscape(location)
	url := fmt.Sprintf(baseURL, c.apiKey, encodedLocation)

	fmt.Printf("WeatherAPI: Buscando clima para cidade: %s\n", city)
	fmt.Printf("WeatherAPI: URL: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("WeatherAPI: Resposta status: %d\n", resp.StatusCode)

	if resp.StatusCode == 400 {
		return nil, fmt.Errorf("cidade não encontrada")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var result model.WeatherAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %v", err)
	}

	return &result, nil
}
