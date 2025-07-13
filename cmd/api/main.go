// @title          Marketplace API
// @version        0.1
// @description    Simple product listing service for GoLang practice.
// @contact.name   Harrison Lisk
// @host           localhost:8080
// @BasePath       /
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/db"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/handlers"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/hariyaki/GoLang-Marketplace-Project/docs"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/cache"
)

func main() {

	//Setup DB
	dsn := os.Getenv("DB_DSN")
	database, err := db.Open(dsn)
	store := listings.NewStore(database)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer database.Close()

	//Setup Cache
	cacheTTL := 30 * time.Second
	redisAddr := os.Getenv("REDIS_ADDR")
	cacher := cache.New(redisAddr, cacheTTL)

	//Set up HTTP handlers
	mux := http.NewServeMux()

	mux.Handle("/healthz", handlers.HealthHandler{})

	// Handle listings request
	// Post = create
	// Get = list / search
	mux.Handle("/listings", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.PostListingHandler{
				Store: store,
			}.ServeHTTP(w, r)

		case http.MethodGet:
			handlers.GetListingsHandler{
				Store: store,
				Cache: cacher,
			}.ServeHTTP(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	//Single Listing GET
	mux.Handle("/listings/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.GetListingHandler{
			Store: store,
			Cache: cacher,
		}.ServeHTTP(w, r)
	}))

	mux.Handle("/docs/", httpSwagger.WrapHandler)

	//Create the server struct
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//A channel used to confirm when ListenAndServe() has returned
	idleConnsClosed := make(chan struct{})

	//Start serving in a goroutine
	//"Graceful Shutdown" Format is from https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
	go func() {
		log.Println("HTTP server starting on", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Any error other than ErrServerClosed is unexpected.
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped accepting new connections.")
		close(idleConnsClosed)
	}()

	//Block until SIGINT or SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Caught signal %s. Shutting down…", sig)

	//Perform graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		// Graceful shutdown didn’t finish in time.
		log.Printf("Shutdown timed out: %v. Forcing close.", err)
		if err := server.Close(); err != nil {
			log.Printf("Forced close failed: %v", err)
		}
	}

	//Wait until ListenAndServe() has actually returned
	<-idleConnsClosed
	log.Println("Shutdown complete.")
}
