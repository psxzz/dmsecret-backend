package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
	"github.com/psxzz/dmsecret-backend/internal/config"
	"github.com/psxzz/dmsecret-backend/internal/repository/secrets"
	"github.com/psxzz/dmsecret-backend/internal/server"
	"github.com/psxzz/dmsecret-backend/internal/server/middlewares"
	"github.com/psxzz/dmsecret-backend/internal/service"
)

const defaultPort = ":3333"

func main() {
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

	secretsRepository, err := secrets.New(cfg.ValkeyConnString)
	if err != nil {
		panic(err)
	}

	svc := service.New(secretsRepository)

	srv := server.NewServer(svc)

	rg := r.Group("/api/v1")
	public.RegisterHandlers(rg, srv)

	log.Fatal(r.Run(defaultPort))
	// log.Fatal(r.RunTLS(defaultPort, "localhost.crt", "localhost.key")) // for local development
}
