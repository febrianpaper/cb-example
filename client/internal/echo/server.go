package echo

import (
	"fbriansyah/client/internal/usecase"
	"fmt"
	"net/http"
)

type Server struct {
	uc      usecase.IUseCases
	address string
}

func New(uc usecase.IUseCases, address string) *Server {
	return &Server{
		uc:      uc,
		address: address,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    s.address,
		Handler: s.Routes(),
	}
	fmt.Printf("server is running on %s\n", s.address)
	return srv.ListenAndServe()
}
