package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/okishiro/pidgey/internal/config"
	one "github.com/okishiro/pidgey/internal/http/handlers/hello"

	//"github.com/okishiro/pidgey/internal/storage"
	"github.com/okishiro/pidgey/internal/storage/sqlite"
)

func main() {
	fmt.Println("hello")
	cfg := config.MustLoad()

	database, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialised")

	router := http.NewServeMux()

	router.HandleFunc("POST /movie/{name}", one.CreateMovie(database))
	//router.HandleFunc("GET /getdata/{id}", one.GetData(storee))
	router.HandleFunc("POST /create/{name}", one.CreateAccount(database, cfg.Storage_path))

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server forced shut %v", err)
	}
	log.Println("server shut")
}
