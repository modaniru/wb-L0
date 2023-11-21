package server

import (
	"html/template"
	"wb-l0/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Server struct{
	storage *storage.OrderStorage
	template *template.Template
}

func NewServer(storage *storage.OrderStorage) *Server{
	template, _ := template.ParseFiles("./template/order.html")
	return &Server{storage: storage, template: template}
}

func (s *Server) InitRouter() *chi.Mux{
	router := chi.NewRouter()
	
	r := router.With(s.Test)
	r.Post("/", s.SaveOrder)
	r.Get("/form", s.GetForm)
	r.Get("/order/{uid}", s.GetOrder)
	return router
}

