package shortener

import (
	"url_shortener/src/repository"

	"github.com/gin-gonic/gin"
)

func SetupHTTPHandler(parentRouter *gin.RouterGroup, svc Service, shortenerRepo repository.URLShortenerRepository) {

	createShortenerHandler := MakeCreateShortenerEndpoint(svc)
	redirectURLHandler := MakeRedirectURLEndpoint(svc)

	r := parentRouter

	r.POST("/shorten", createShortenerHandler)
	r.GET("/:code", redirectURLHandler)
}
