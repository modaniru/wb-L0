package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"wb-l0/internal/entity"

	"github.com/go-chi/chi/v5"
)

func (s *Server) SaveOrder(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.Body)
	log.Println("read r.Body")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("parsing body")
	order := entity.Order{}
	err = json.Unmarshal(body, &order)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("push to db")
	err = s.storage.SaveOrder(r.Context(), order.OrderUid, body)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request){
	uid := chi.URLParam(r, "uid")
	data, err := s.storage.GetByUid(r.Context(), uid)
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var buf bytes.Buffer
	json.Indent(&buf, data, "", "\t")
	w.Write(buf.Bytes())
}