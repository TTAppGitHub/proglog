package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type httpServer struct {
	Log *Log
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func NewHTTPServer(address string) *http.Server {
	httpserv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpserv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpserv.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    address,
		Handler: r,
	}

}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
