package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    
    "github.com/suraif16/webpage-analyzer/internal/core/services"
    "github.com/suraif16/webpage-analyzer/internal/handlers"
    httpClient "github.com/suraif16/webpage-analyzer/internal/infrastructure/http/client"
    "github.com/suraif16/webpage-analyzer/internal/infrastructure/parser"
    "github.com/suraif16/webpage-analyzer/internal/middleware"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/suraif16/webpage-analyzer/docs"
)
func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Initialize dependencies
    httpClient := httpClient.NewHTTPClient(10 * time.Second)
    htmlParser := parser.NewHTMLParser()
    analyzerService := services.NewAnalyzerService(httpClient, htmlParser, logger.Sugar())
    analyzerHandler := handlers.NewAnalyzerHandler(analyzerService)

    // Setup Gin
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.Logger(logger))
    r.Use(middleware.CORS())

    // Swagger route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Routes
    r.POST("/analyze", analyzerHandler.Analyze)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Server configuration
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    // Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("failed to start server", zap.Error(err))
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("server forced to shutdown:", zap.Error(err))
    }

    logger.Info("server exited properly")
}