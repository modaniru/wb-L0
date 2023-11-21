package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	natssub "wb-l0/internal/nats-sub"
	"wb-l0/internal/server"
	"wb-l0/internal/storage"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
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

	subscriber := natssub.NewSubscriber(orderStorage)
	
	conn, err := stan.Connect("prod", "subscriber")
	if err != nil{
		log.Fatal(err.Error())
	}
	defer conn.Close()

	conn.Subscribe("test", subscriber.GetMsgHandler(), stan.DeliverAllAvailable())
	
	// wait nats

    http.ListenAndServe(":80", router)
}

func initLogger(){
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	slog.Info("slog was init...")
}