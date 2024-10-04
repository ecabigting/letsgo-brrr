package server

import "net/http"

type Server struct {
	Mux                     http.ServeMux
	SubscriberMessageBuffer int
}

func (s *Server) NewServer() {
	s.SubscriberMessageBuffer = 10
	s.Mux.Handle("/", http.FileServer(http.Dir("../../htmx")))
}
