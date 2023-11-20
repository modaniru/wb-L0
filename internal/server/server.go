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
	router.Post("/", s.SaveOrder)
	router.Get("/{uid}", s.GetOrder)
	return router
}