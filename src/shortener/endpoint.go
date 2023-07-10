package shortener

import (
	"net/http"
	"url_shortener/src/shared/model"

	"github.com/gin-gonic/gin"
)

func MakeCreateShortenerEndpoint(svc Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateShortenerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		inputData := model.ShortURL{
			ID:          req.ID,
			FullURL:     req.FullURL,
			ShortCode:   req.ShortCode,
			AccessCount: 0,
		}

		status, err := svc.CreateShortener(c, &inputData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := CreateShortenerResponse{
			Status: status,
		}

		c.JSON(200, resp)
	}
}

type CreateShortenerRequest struct {
	ID          string `json:"id"`
	FullURL     string `json:"fullURL"`
	ShortCode   string `json:"shortCode"`
	AccessCount int    `json:"accessCount"`
}

type CreateShortenerResponse struct {
	Status string `json:"status"`
}

func MakeRedirectURLEndpoint(svc Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		url, err := svc.RedirectURL(c, c.Param("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := RedirectURLResponse{
			FullURL: url,
		}

		c.JSON(200, resp)
	}
}

type RedirectURLResponse struct {
	FullURL string `json:"fullUrl"`
}
