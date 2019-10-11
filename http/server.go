package http

import (
	"net"
	"net/http"
	"net/url"

	disbursement "github.com/jfpalngipang/fund-disbursement"
	"github.com/pressly/chi"
)

type Server struct {
	ln net.Listener

	// Services
	DisbursementService disbursement.DisbursementService
	// Server options
	Addr string
	// Host string
}

func NewServer() *Server {
	return &Server{}
}

// Open opens the server.
func (s *Server) Open() error {
	// Open listener on specified bind address.
	// Use HTTPS port if autocert is enabled.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	// Start HTTP server.
	http.Serve(s.ln, s.router())

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Mount("/disbursement", s.disbursementHandler())
		r.Get("/health", healthHandler)
		r.Get("/readiness", readinessHandler)
	})

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) disbursementHandler() *disbursementHandler {
	h := newDisbursementHandler()
	h.baseUrl = s.URL()
	h.disbursementService = s.DisbursementService
	return h

}

func (s *Server) URL() url.URL {
	if s.ln == nil {
		return url.URL{}
	}

	return url.URL{Scheme: "http", Host: s.ln.Addr().String()}
}
