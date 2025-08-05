package viacep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oliveiracmorais/cep-clima/model"
)

const baseURL = "https://viacep.com.br/ws/%s/json/"

type Client interface {
	GetAddress(cep string) (*model.ViaCEPResponse, error)
}

type ViaCEPClient struct{}

func NewViaCEPClient() Client {
	return &ViaCEPClient{}
}

func (c *ViaCEPClient) GetAddress(cep string) (*model.ViaCEPResponse, error) {
	url := fmt.Sprintf(baseURL, cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return nil, fmt.Errorf("invalid zipcode")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var result model.ViaCEPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta: %v", err)
	}

	if result.Erro == "true" {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &result, nil
}
