package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	
	"json-crud-service/internal/infrastructure/repository"
	"json-crud-service/internal/presentation/handler"
	"json-crud-service/internal/presentation/middleware"
	"json-crud-service/internal/usecase"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Initialize dependencies
	jsonRepo := repository.NewMemoryJSONRepository()
	jsonUsecase := usecase.NewJSONUsecase(jsonRepo)
	jsonHandler := handler.NewJSONHandler(jsonUsecase)
	
	// Setup HTTP server
	mux := http.NewServeMux()
	jsonHandler.SetupRoutes(mux)
	
	// Apply middleware
	handler := middleware.CORS(middleware.Logging(mux))
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	
	// Start server in a goroutine
	go func() {
		fmt.Printf("Server starting on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("Server shutting down...")
	
	// Here you could add graceful shutdown logic if needed
	// For now, we'll just exit
	fmt.Println("Server stopped")
}
