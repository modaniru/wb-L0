package server

import (
	"bytes"
	"encoding/json"

	"io"
	"log"
	"net/http"

	"github.com/modaniru/wb-L0/internal/entity"
)

func (s *Server) SaveOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	log.Println("read r.Body")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("parsing body")
	order := entity.Order{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("push to db")
	err = s.storage.SaveOrder(r.Context(), order.OrderUid, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	data, err := s.storage.GetByUid(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var buf bytes.Buffer
	json.Indent(&buf, data, "", "\t")
	w.Write(buf.Bytes())
}

func (s *Server) GetForm(w http.ResponseWriter, r *http.Request) {
	model := map[string]any{}
	uid := r.URL.Query().Get("uid")

	rc, err := s.storage.GetRowsCount(r.Context())
	if err != nil{
		model["error"] = err.Error()
		s.template.Execute(w, model)
		return
	}

	model["cache"] = rc.Cache
	model["db"] = rc.DB
	if uid == "" {
		s.template.Execute(w, model)
		return
	}

	data, err := s.storage.GetByUid(r.Context(), uid)
	if err != nil {
		model["error"] = err.Error()
		s.template.Execute(w, model)
		return
	}

	model["order"] = string(data)
	s.template.Execute(w, model)
}
