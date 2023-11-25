package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	natssub "github.com/modaniru/wb-L0/internal/nats-sub"
	"github.com/modaniru/wb-L0/internal/server"
	"github.com/modaniru/wb-L0/internal/storage"
	"github.com/nats-io/stan.go"
)

func App() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5555?sslmode=disable")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	orderStorage := storage.NewOrderStorage(db, storage.NewInmemoryCache())
	server := server.NewServer(orderStorage)
	router := server.InitRouter()

	subscriber := natssub.NewSubscriber(orderStorage)

	conn, err := stan.Connect("prod", "subscriber")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	_, err = conn.Subscribe("test", subscriber.GetMsgHandler(), stan.DurableName("client"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	httpServer := http.Server{
		Addr:    ":80",
		Handler: router,
	}
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	<-signalChan

	shutdown := context.Background()

	if err := httpServer.Shutdown(shutdown); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("finished")
}

func initLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	slog.Info("slog was init...")
}
