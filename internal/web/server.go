package web

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type templateContainer interface {
	Get(string) (*Template, error)
}

type Server struct {
	container templateContainer
	address   string
}

func NewServer(address string, container templateContainer) *Server {
	return &Server{
		container: container,
		address:   address,
	}
}

func (s *Server) handlers() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page := vars["page"]

		tpl, err := s.container.Get(page)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		if err = tpl.Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	return router
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	server := http.Server{
		Addr:    s.address,
		Handler: s.handlers(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	log.Println("web server started")
	<-ctx.Done()
	log.Println("shutting down the web server")

	ctxShutDown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
	}()

	return server.Shutdown(ctxShutDown)
}
