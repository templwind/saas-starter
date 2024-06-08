package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/templwind/sass-starter/internal/config"
	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/modules"

	_ "github.com/a-h/templ"
	"github.com/joho/godotenv"
)

//go:embed etc/config.yaml
var configFile embed.FS

//go:embed db/migrations/*.sql
var embeddedMigrations embed.FS

func main() {
	err := godotenv.Load(".env", ".envrc")
	if err != nil {
		log.Println("Error loading .env file")
	}

	configBytes, err := configFile.ReadFile("etc/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Expand environment variables
	configBytes = []byte(os.ExpandEnv(string(configBytes)))

	var c config.Config
	err = config.LoadConfigFromYamlBytes(configBytes, &c)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set the embedded migrations
	c.EmbeddedMigrations = embeddedMigrations

	// Create a context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create the service context instance
	svcCtx := svc.NewServiceContext(ctx, &c)

	// Set up the HTTP server
	mux := http.NewServeMux()

	// Static file handler
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("assets")))
	mux.Handle("/static/", staticHandler)

	// Register modules (assuming RegisterAll accepts ServeMux instead of Echo)
	modules.RegisterAll(svcCtx, mux)

	// Middleware: Remove trailing slash and redirect
	mux.HandleFunc("/", middleware.RemoveTrailingSlash(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})))

	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	// Start the server in a new goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8888: %v\n", err)
		}
	}()

	// Wait for the context to be canceled (which happens when we receive an interrupt signal)
	<-ctx.Done()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	// Create a new context for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
