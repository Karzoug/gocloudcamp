package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Karzoug/gocloudcamp/internal/player"
	"github.com/Karzoug/gocloudcamp/internal/playlist/memory"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	pl := memory.New()
	player := player.New(pl)
	defer player.Close()

	<-ctx.Done()
	log.Println("stop the server gracefully")
}
