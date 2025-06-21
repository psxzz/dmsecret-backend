package middlewares

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
)

func WithOAPIRequestValidation(oapiPath string) gin.HandlerFunc {
	path, err := filepath.Abs(oapiPath)
	if err != nil {
		panic(err)
	}

	mw, err := ginmiddleware.OapiValidatorFromYamlFile(path)
	if err != nil {
		panic(err)
	}

	return mw
}
