package server

import (
	"wb-l0/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Server struct{
	storage *storage.OrderStorage
}

func NewServer(storage *storage.OrderStorage) *Server{
	return &Server{storage: storage}
}

func (s *Server) InitRouter() *chi.Mux{
	router := chi.NewRouter()
	r := router.With(s.Test)
	r.Post("/", s.SaveOrder)
	r.Get("/{uid}", s.GetOrder)
	return router
}