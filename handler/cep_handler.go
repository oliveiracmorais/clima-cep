package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/oliveiracmorais/cep-clima/model"
	"github.com/oliveiracmorais/cep-clima/service"
)

type CEPHandler struct {
	Service service.CEPClimaServiceInterface
}

func NewCEPHandler() *CEPHandler {
	return &CEPHandler{
		Service: service.NewCEPClimaService(),
	}
}

func (h *CEPHandler) GetClimaPorCEP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/cep/")
	if path == r.URL.Path {
		http.Error(w, "Endpoint não encontrado", http.StatusNotFound)
		return
	}

	cep := path

	fmt.Printf("Handler: Recebida requisição para CEP: %s\n", cep)

	response, err := h.Service.GetClimaPorCEP(cep)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if err.Error() == "invalid zipcode" {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "http 422 - invalid zipcode"})
		} else if err.Error() == "can not find zipcode" {
			w.WriteHeader(http.StatusNotFound) // 404
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "http 404 - can not find zipcode"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "erro interno"})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *CEPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
