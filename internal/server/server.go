package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
)

//go:generate go tool oapi-codegen -config ../../api/public/cfg.yaml ../../api/public/api.yaml
var _ public.ServerInterface = (*Server)(nil)

type Repository interface {
	CreateSecret(ctx context.Context, payload string) (string, error)
}

type Server struct {
	repo Repository
}

func NewServer(repo Repository) *Server {
	s := &Server{
		repo: repo,
	}
	return s
}

func (s *Server) GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, public.AliveResponse{Text: "alive"})
}
