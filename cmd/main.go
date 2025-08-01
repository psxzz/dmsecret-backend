package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
	"github.com/psxzz/dmsecret-backend/internal/config"
	"github.com/psxzz/dmsecret-backend/internal/cryptographer"
	"github.com/psxzz/dmsecret-backend/internal/repository/secrets"
	"github.com/psxzz/dmsecret-backend/internal/server"
	"github.com/psxzz/dmsecret-backend/internal/server/middlewares"
	"github.com/psxzz/dmsecret-backend/internal/service"
)

const defaultPort = ":3333"

func main() {
	log.Fatal(run())
}

func run() error {
	ctx := context.Background()
	_ = ctx

	cfg, err := config.Overload()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		middlewares.WithCORSCheck(),
		middlewares.WithOAPIRequestValidation(cfg.OAPIPath),
	)

	crypt, err := cryptographer.New(cfg.CryptoKey)
	if err != nil {
		panic(err)
	}

	secretsRepository, err := secrets.New(cfg.ValkeyConnString, crypt)
	if err != nil {
		panic(err)
	}
	defer secretsRepository.Close()

	svc := service.New(secretsRepository)

	srv := server.NewServer(svc)

	rg := r.Group("/api/v1")
	public.RegisterHandlers(rg, srv)

	err = r.Run(defaultPort)
	// err = r.RunTLS(r.RunTLS(defaultPort, "ssl-cert-snakeoil.pem", "ssl-cert-snakeoil.key")) // for local development
	return fmt.Errorf("server error: %w", err)
}
