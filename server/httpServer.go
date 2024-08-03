package server

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	SubscriberMsgBuffer int
	Mux                 http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		SubscriberMsgBuffer: 10,
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get CWD setting to relative './' : %s\n", err)
		cwd = "./"
	}

	htmxDir := fmt.Sprintf("%s/htmx", cwd)

	s.Mux.Handle("/", http.FileServer(http.Dir(htmxDir)))
	return s
}
