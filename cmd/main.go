package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
	"github.com/psxzz/dmsecret-backend/internal/config"
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
	_ = cfg
	srv := server.NewServer()

	r := gin.New()
	r.Use(server.CORSMiddleware())
	rg := r.Group("/api/v1")

	public.RegisterHandlers(rg, srv)

	log.Fatal(r.Run(defaultPort))
	// log.Fatal(r.RunTLS(defaultPort, "localhost.crt", "localhost.key")) // for local development
}
