package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	viaCepBaseUrl       string = "http://viacep.com.br/ws/"
	viaCepFormatUrl     string = "/json/"
	brasilApiCepBaseUrl        = "https://brasilapi.com.br/api/cep/v1/"
	timeout                    = time.Second
)

type apiResponse struct {
	api     string `json:"api"`
	message string `json:"message"`
}

func main() {

	cep := "01153000"

	resCha := make(chan *apiResponse, 1)

	go fetchViaCep(cep, resCha)
	go fetchBrasilApiCep(cep, resCha)

	// Seleciona o resultado mais rápido
	select {
	case result := <-resCha:
		if result != nil {
			fmt.Println("Resultado da api: ", result.api, "Response:", result.message)
		} else {
			fmt.Println("Erro ao obter resultado da brasilapi.")
		}
	case <-time.After(timeout):
		fmt.Println("Timeout: Nenhuma resposta dentro do tempo limite.")
	}
}

func fetchViaCep(cep string, res chan<- *apiResponse) {

	response, err := http.Get(viaCepBaseUrl + cep + viaCepFormatUrl)

	if err != nil {
		res <- nil
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		res <- nil
		return
	}

	result := apiResponse{api: "ViaCep", message: string(body)}
	res <- &result
}

func fetchBrasilApiCep(cep string, res chan<- *apiResponse) {
	//Como a api do BrasilAPi só retorna erro, deixei o sleep para simular a Via cep respondendo primeiro
	//time.Sleep(time.Second * 100)
	response, err := http.Get(brasilApiCepBaseUrl + cep)

	if err != nil {
		res <- nil
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		res <- nil
		return
	}

	result := apiResponse{api: "BrasilApi", message: string(body)}
	res <- &result
}
