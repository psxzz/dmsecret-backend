package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/psxzz/dmsecret-backend/api/public"
)

func (s *Server) PostSecrets(c *gin.Context) {
	var body public.PostSecretsJSONRequestBody

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, public.ValidationError{Error: err.Error()})
		return
	}

	secretID, err := s.svc.CreateSecret(c.Request.Context(), body.Payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, public.SecretsOut{SecretID: secretID})
}
