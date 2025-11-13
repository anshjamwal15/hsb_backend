package server

import (
	"context"
	"log"
	"net/http"

	"github.com/anshjamwal15/hsb_backend/internal/config"
	httpPkg "github.com/anshjamwal15/hsb_backend/internal/http"
	"github.com/anshjamwal15/hsb_backend/internal/infrastructure/database"
)

type Server struct {
	httpServer *http.Server
	db         *database.MongoDB
}

func NewServer(cfg *config.Config, db *database.MongoDB) *Server {
	// Setup router with all routes
	router := httpPkg.SetupRouter(db, cfg)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &Server{
		httpServer: httpServer,
		db:         db,
	}
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
