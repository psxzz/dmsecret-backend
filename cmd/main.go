package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
	"github.com/psxzz/dmsecret-backend/internal/config"
	"github.com/psxzz/dmsecret-backend/internal/database/valkey"
	"github.com/psxzz/dmsecret-backend/internal/repository/secrets"
	"github.com/psxzz/dmsecret-backend/internal/server"
	"github.com/psxzz/dmsecret-backend/internal/server/middlewares"
	"github.com/psxzz/dmsecret-backend/internal/service"
)

const defaultPort = ":3333"

func main() {
	ctx := context.Background()
	_ = ctx

	cfg, err := config.Create()
	if err != nil {
		panic(err)
	}

	keyValueDB, err := valkey.New(cfg.ValkeyConnString)
	if err != nil {
		panic(err)
	}

	// postgresDB, err := postgres.New(ctx, cfg.PGConnString)
	// if err != nil {
	// 	panic(err)
	// }
	// _ = postgresDB

	r := gin.New()
	r.Use(
		middlewares.WithCORSCheck(),
		middlewares.WithOAPIRequestValidation("./api/public/api.yaml"),
		gin.Logger(),
		gin.Recovery(),
	)

	secretsRepository := secrets.New(keyValueDB)

	repo := service.New(secretsRepository)

	srv := server.NewServer(repo)

	rg := r.Group("/api/v1")
	public.RegisterHandlers(rg, srv)

	log.Fatal(r.Run(defaultPort))
	// log.Fatal(r.RunTLS(defaultPort, "localhost.crt", "localhost.key")) // for local development
}
