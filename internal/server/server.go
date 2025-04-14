package server

import (
	"context"
	"wallet-api-service/internal/api"
	"wallet-api-service/internal/db"
	"wallet-api-service/internal/db/memdb"
	"wallet-api-service/internal/config"
	"wallet-api-service/internal/kafka"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg     *config.Config
	db      db.DB
	api     *api.API
	kafka   *kafka.Client
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{cfg: cfg}

	s.db = memdb.New()
	s.kafka = kafka.New(*s.cfg)
	s.api = api.New(*s.cfg, s.db, s.kafka)

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.api.Serve(ctx)
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Info().Msg("Graceful server shutdown")
	
	if err := s.api.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down HTTP server")
	}

	if s.kafka != nil {
		if err := s.kafka.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Kafka connection")
		}
	}
}
