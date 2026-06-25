package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FerT99/gestor-terrenos-service/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database.ConnectDB()
	var parcelaID string
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM parcelas WHERE nombre = 'Parcela Principal' LIMIT 1").Scan(&parcelaID)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/abonos", nil)
	req.Header.Set("X-Parcela-Id", parcelaID)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
