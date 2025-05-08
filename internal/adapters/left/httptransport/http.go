package http

import (
	ports "1337b04rd/internal/ports/left"
	"net/http"
)

type Server struct {
	addr   string
	router *http.ServeMux
	svc    ports.APIPort
}

func NewHTTPServer(svc ports.APIPort) *Server {
	router := newRouter(svc)

	addr := ":8080"
	return &Server{
		addr:   addr,
		router: router,
		svc:    svc,
	}
}

func newRouter(svc ports.APIPort) *http.ServeMux {
	router := http.NewServeMux()

	RegisterRoutes(svc, router)
	return router
}

func (s *Server) Serve() error {
	wrappedRouter := Chain(s.router, WithSession(s.svc))

	srv := &http.Server{
		Addr:    s.addr,
		Handler: wrappedRouter,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout:  120 * time.Second,
	}

	return srv.ListenAndServe()
}
