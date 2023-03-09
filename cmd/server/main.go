package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Karzoug/gocloudcamp/internal/config"
	"github.com/Karzoug/gocloudcamp/internal/player"
	"github.com/Karzoug/gocloudcamp/internal/playlist"
	"github.com/Karzoug/gocloudcamp/internal/playlist/file"
	"github.com/Karzoug/gocloudcamp/internal/playlist/memory"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.New()
	if err := cfg.Load(); err != nil {
		log.Fatalf("load config error: %v", err)
	}

	var pl playlist.Playlist
	if cfg.IsStoreInMemory() {
		pl = memory.New()
	} else {
		var err error
		pl, err = file.New(cfg)
		if err != nil {
			log.Fatalf("create playlist error: %v", err)
		}
	}

	player := player.New(pl)
	defer player.Close()

	<-ctx.Done()
	log.Println("stop the server gracefully")
}
