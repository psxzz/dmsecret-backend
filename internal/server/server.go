package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/psxzz/dmsecret-backend/api/public"
)

//go:generate go tool oapi-codegen -config ../../api/public/cfg.yaml ../../api/public/api.yaml
var _ public.ServerInterface = (*Server)(nil)

type Service interface {
	CreateSecret(ctx context.Context, payload string) (string, error)
	GetSecretByID(ctx context.Context, id uuid.UUID) (*string, error)
}

type Server struct {
	svc Service
}

func NewServer(svc Service) *Server {
	s := &Server{
		svc: svc,
	}
	return s
}

func (s *Server) GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, public.AliveResponse{Text: "alive"})
}
