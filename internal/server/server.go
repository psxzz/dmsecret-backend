package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
)

//go:generate go tool oapi-codegen -config ../../api/public/cfg.yaml ../../api/public/api.yaml
var _ public.ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, public.AliveResponse{Text: "alive"})
}
