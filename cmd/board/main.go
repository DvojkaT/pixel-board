package main

import (
	"board/internal/canvas"
	"board/internal/handler"
	"board/internal/hub"
	"board/internal/ratelimit"
	"board/static"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gameCanvas := canvas.NewCanvas(100, 100)
	gameLimiter := ratelimit.NewLimiter()
	gameHub, err := hub.NewHub(gameCanvas, gameLimiter)
	if err != nil {
		log.Fatal(err)
	}

	gameHandler := handler.NewHandler(gameHub, static.FS)
	gameMux := http.NewServeMux()
	gameHandler.RegisterRoutes(gameMux)

	go func() {
		if err := http.ListenAndServe(":8080", gameMux); err != nil {
			log.Fatal(err)
		}
	}()

	if err := gameHub.Run(); err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	log.Println("server started on :8080")
	<-sigCh
	log.Println("shutting down")
}
