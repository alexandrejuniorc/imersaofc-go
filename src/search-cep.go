package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

type BrasilAPIResponse struct {
	CEP          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type UnifiedResponse struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	UF         string `json:"uf"`
	Source     string `json:"source"`
}

func fetchViaCEP(context context.Context, cep string, resultChan chan<- UnifiedResponse) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	request, err := http.NewRequestWithContext(context, "GET", url, nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	var viaCEPResponse ViaCEPResponse
	if err := json.NewDecoder(response.Body).Decode(&viaCEPResponse); err != nil {
		return
	}

	resultChan <- UnifiedResponse{
		CEP:        viaCEPResponse.CEP,
		Logradouro: viaCEPResponse.Logradouro,
		Bairro:     viaCEPResponse.Bairro,
		Cidade:     viaCEPResponse.Localidade,
		UF:         viaCEPResponse.UF,
		Source:     "ViaCEP",
	}
}

func fetchBrasilAPI(context context.Context, cep string, resultChan chan<- UnifiedResponse) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	request, err := http.NewRequestWithContext(context, "GET", url, nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	var brasilAPIResponse BrasilAPIResponse
	if err := json.NewDecoder(response.Body).Decode(&brasilAPIResponse); err != nil {
		return
	}

	resultChan <- UnifiedResponse{
		CEP:        brasilAPIResponse.CEP,
		Logradouro: brasilAPIResponse.Street,
		Bairro:     brasilAPIResponse.Neighborhood,
		Cidade:     brasilAPIResponse.City,
		UF:         brasilAPIResponse.State,
		Source:     "BrasilAPI",
	}
}

func searchCEP(writer http.ResponseWriter, request *http.Request) {
	cep := request.URL.Query().Get("cep")
	if cep == "" {
		http.Error(writer, "CEP is required", http.StatusBadRequest)
		return
	}

	// Create a new context with a timeout of 5 seconds
	context, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel() // Cancel the context to release resources

	resultChan := make(chan UnifiedResponse, 2) // Buffered channel with capacity of 2

	// Start goroutines to fetch data from ViaCEP and BrasilAPI
	go fetchViaCEP(context, cep, resultChan)
	go fetchBrasilAPI(context, cep, resultChan)

	// Wait for the first response
	select {
	case result := <-resultChan:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(result)
		return

	case <-context.Done():
		http.Error(writer, "Timeout exceeded", http.StatusRequestTimeout)
	}
}
