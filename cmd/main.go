package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
	"github.com/psxzz/dmsecret-backend/internal/config"
	"github.com/psxzz/dmsecret-backend/internal/database/key_value"
	"github.com/psxzz/dmsecret-backend/internal/database/postgres"
	"github.com/psxzz/dmsecret-backend/internal/server"
)

const defaultPort = ":3333"

func main() {
	ctx := context.Background()
	_ = ctx

	cfg, err := config.Create()
	if err != nil {
		panic(err)
	}

	keyValueDB, err := key_value.New(cfg)
	if err != nil {
		panic(err)
	}
	_ = keyValueDB

	postgresDB, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}
	_ = postgresDB

	srv := server.NewServer()

	r := gin.New()
	r.Use(server.CORSMiddleware())
	rg := r.Group("/api/v1")

	public.RegisterHandlers(rg, srv)

	log.Fatal(r.Run(defaultPort))
	// log.Fatal(r.RunTLS(defaultPort, "localhost.crt", "localhost.key")) // for local development
}
