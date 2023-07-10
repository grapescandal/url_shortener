package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_shortener/src/config"
	"url_shortener/src/repository"
	shortener "url_shortener/src/shortener"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
)

func Run() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		panic(errors.Wrap(err, "failed to load config from env"))
	}

	fmt.Println("Load config success")

	var urlShortenerRepo repository.URLShortenerRepository
	var shortenerSvc shortener.Service

	urlShortenerRepo_ := repository.NewURLShortenerRepositoryImpl()
	urlShortenerRepo = urlShortenerRepo_

	shortenerSvc_ := shortener.NewServiceImpl(urlShortenerRepo)
	shortenerSvc = shortenerSvc_

	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	apiV1Router := router.Group("/api/v1")
	shortener.SetupHTTPHandler(apiV1Router, shortenerSvc, urlShortenerRepo)

	router.Use(HandleCORS())

	server := &http.Server{
		Addr:    appConfig.HTTPHost,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("failed to start server, err: %s\n", err.Error())
		}
	}()

	// do graceful stop...
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	fmt.Println("subscribed for SIGTERM, SIGINT signal for graceful stop")

	<-gracefulStop
	fmt.Println("gracefully stop...")

	// clean up other resources...
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server with err: %s", err.Error())
	}

	fmt.Println("gracefully stop... done")
}

func HandleCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
