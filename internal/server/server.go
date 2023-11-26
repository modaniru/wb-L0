package server

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/modaniru/wb-L0/internal/storage"
)

type Server struct {
	storage  *storage.OrderStorage
	template *template.Template
}

func NewServer(storage *storage.OrderStorage) *Server {
	template, _ := template.ParseFiles("./static/html/main.html")
	return &Server{storage: storage, template: template}
}

func (s *Server) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	fs := http.FileServer(http.Dir("static"))

	r := router.With(s.Test)
	r.Post("/", s.SaveOrder)
	r.Get("/form", s.GetForm)
	r.Get("/order", s.GetOrder)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	return router
}
