package cmd

import (
	"GhostyLink/api/public"
	"GhostyLink/internal/config"
	"GhostyLink/internal/server"
	"context"
	"github.com/gin-gonic/gin"
	"log"
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

	//log.Fatal(r.Run(defaultPort))
	log.Fatal(r.RunTLS(defaultPort, "localhost.crt", "localhost.key")) // for local development
}
