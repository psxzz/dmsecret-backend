package server

import (
	"GhostyLink/api/public"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	s := &Server{}
	return s
}

func (s *Server) GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, public.AliveResponse{Text: "alive"})
}
