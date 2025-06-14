package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
)

type Server struct{}

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, public.AliveResponse{Text: "alive"})
}
