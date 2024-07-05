package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/wopoczynski/todoapp/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = godotenv.Load()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Server error") //nolint:forbidigo
	}
}
