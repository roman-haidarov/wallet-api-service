package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-api-service/internal/config"
	"wallet-api-service/internal/logger"
	"wallet-api-service/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

func main() {
	logger.Init(zerolog.DebugLevel)
	configPath := pflag.StringP("config", "c", "", "path to config file")
	pflag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(context.Background()); err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	}()

	log.Info().Msg("Server started")

	<-sigChan
	log.Info().Msg("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info().Msg("Server stopped")
}
