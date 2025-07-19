package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/psxzz/dmsecret-backend/api/public"
)

func (s *Server) GetSecrets(c *gin.Context, params public.GetSecretsParams) {
	id, err := uuid.Parse(params.SecretID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, public.ValidationError{Error: err.Error()})
		return
	}

	payload, err := s.svc.GetSecretByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if payload == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "secret not found"})
		return
	}

	c.JSON(http.StatusOK, public.GetSecretsOut{Payload: *payload})
}
