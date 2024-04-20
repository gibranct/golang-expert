package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type Dollar struct {
	Price string `json:"price"`
}

func saveToFile(dollar *Dollar) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path.Join(cwd, "cotacao.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("Dólar:{" + dollar.Price + "}\n")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	select {
	case <-ctxWithTimeout.Done():
		log.Fatal("Failed to get dollar price in an decent time")
	case <-time.After(time.Millisecond * 299):
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var dollar Dollar
		json.Unmarshal(body, &dollar)
		err = saveToFile(&dollar)
		if err != nil {
			log.Fatal("Failed to save to file: " + err.Error())
		}
	}
}
