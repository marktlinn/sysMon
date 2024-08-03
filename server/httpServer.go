package server

import "net/http"

type Server struct {
	SubscriberMsgBuffer int
	Mux                 http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		SubscriberMsgBuffer: 10,
	}

	s.Mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	return s
}
