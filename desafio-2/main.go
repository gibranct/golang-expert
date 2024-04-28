package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CepAPIResponse struct {
	API  string
	Data string
}

func main() {
	args := os.Args
	cep := "01153000"
	if len(args) > 1 {
		cep = args[1]
	}
	ch := make(chan CepAPIResponse)
	go requestTwo(cep, ch)
	go requestOne(cep, ch)

	select {
	case response := <-ch:
		fmt.Println(response)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}

func requestOne(cep string, ch chan CepAPIResponse) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep

	http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("[BrasilAPI]: Fail to get CEP, err: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("[BrasilAPI]: Fail to get CEP, err: " + err.Error())
	}
	ch <- CepAPIResponse{
		API:  url,
		Data: string(body),
	}
}

func requestTwo(cep string, ch chan CepAPIResponse) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("[ViaCepAPI]: Fail to get CEP, err: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("[ViaCepAPI]: Fail to get CEP, err: " + err.Error())
	}
	ch <- CepAPIResponse{
		API:  url,
		Data: string(body),
	}
}
