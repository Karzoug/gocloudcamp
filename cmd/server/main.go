package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/Karzoug/gocloudcamp/internal/config"
	"github.com/Karzoug/gocloudcamp/internal/grpcapi"
	"github.com/Karzoug/gocloudcamp/internal/player"
	"github.com/Karzoug/gocloudcamp/internal/playlist"
	"github.com/Karzoug/gocloudcamp/internal/playlist/file"
	"github.com/Karzoug/gocloudcamp/internal/playlist/memory"
	"github.com/Karzoug/gocloudcamp/internal/server"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.New()
	if err := cfg.Load(); err != nil {
		log.Fatalf("load config error: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

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

	grpcapi.RegisterPlayerServiceServer(s, server.New(player))

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("stop the server gracefully")
	s.GracefulStop()
}
