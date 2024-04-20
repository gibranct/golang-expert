package main

import (
	"context"
	"database/sql"
	_ "database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com.br/gibranct/desafio/server/clients"
	_ "github.com/mattn/go-sqlite3"
)

type PriceResponse struct {
	Price string `json:"price"`
}

var sqliteDB *sql.DB

func createDBConn() {
	db, err := sql.Open("sqlite3", "./quotation.db")
	if err != nil {
		panic(err)
	}
	sqliteDB = db
	_, err = db.Exec("create table if not exists quotation (dollar jsonb)")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", FetchDollarPriceHandler)
	createDBConn()
	http.ListenAndServe(":8080", mux)
}

func FetchDollarPriceHandler(w http.ResponseWriter, r *http.Request) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	data, err := clients.DollarQuotation(ctxWithTimeout)
	select {
	case <-ctxWithTimeout.Done():
		fmt.Fprintf(os.Stderr, "Request canceled\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	default:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fail to unmarshal body %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = saveData(*data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fail to save body %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(PriceResponse{Price: data.Bid})
	}

}

func saveData(data clients.Quotation) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancel()
	stmt, err := sqliteDB.PrepareContext(ctx, "insert into quotation(dollar) values(?)")
	if err != nil {
		return err
	}
	jsonData, _ := json.Marshal(data)
	_, err = stmt.Exec(jsonData)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var dq []byte
	err = sqliteDB.QueryRow("select dollar->>'USDBRL' from quotation").Scan(&dq)
	if err != nil {
		return err
	}
	return err
}
