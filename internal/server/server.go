package server

import (
	"html/template"

	"github.com/go-chi/chi/v5"
	"github.com/modaniru/wb-L0/internal/storage"
)

type Server struct {
	storage  *storage.OrderStorage
	template *template.Template
}

func NewServer(storage *storage.OrderStorage) *Server {
	template, _ := template.ParseFiles("./template/order.html")
	return &Server{storage: storage, template: template}
}

func (s *Server) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	r := router.With(s.Test)
	r.Post("/", s.SaveOrder)
	r.Get("/form", s.GetForm)
	r.Get("/order", s.GetOrder)

	return router
}
