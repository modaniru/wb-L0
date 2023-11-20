package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"wb-l0/internal/server"
	"wb-l0/internal/storage"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5555?sslmode=disable")
	if err != nil{
		log.Fatal(err)
	}

	err = db.PingContext(context.Background())
	if err != nil{
		log.Fatal(err)
	}
	
	orderStorage := storage.NewOrderStorage(db, storage.NewCache())
	server := server.NewServer(orderStorage)
	router := server.InitRouter()

    http.ListenAndServe(":80", router)
}
