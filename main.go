package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	brasilApiUrl = "https://brasilapi.com.br/api/cep/v1/%s"
	viaCepApiUrl = "http://viacep.com.br/ws/%s/json/"
)

type ApiResponse struct {
	Api     string
	Address map[string]string
}

func main() {
	cep := "01153000"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	apiResponseCh := make(chan ApiResponse, 1)

	go fetchAddressFromAPI(fmt.Sprintf(brasilApiUrl, cep), apiResponseCh)
	go fetchAddressFromAPI(fmt.Sprintf(viaCepApiUrl, cep), apiResponseCh)

	select {
	case result := <-apiResponseCh:
		fmt.Println("API:", result.Api)
		fmt.Println("Address:", result.Address)
	case <-ctx.Done():
		log.Fatal("Timeout")
	}
}

func fetchAddressFromAPI(url string, ch chan ApiResponse) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	client := http.Client{
		Timeout: time.Second * 1,
	}

	r, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		log.Fatalf("Error: API request failed with status code %d", r.StatusCode)
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(data, &apiResponse.Address)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	apiResponse.Api = r.Request.URL.String()

	ch <- apiResponse
}
