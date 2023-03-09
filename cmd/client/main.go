package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/Karzoug/gocloudcamp/internal/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := grpcapi.NewPlayerServiceClient(conn)

	audios := []grpcapi.Audio{
		{
			Name:     "Simple ringtone",
			Duration: durationpb.New(10 * time.Second),
		},
		{
			Name:     "Meowing cat",
			Duration: durationpb.New(18 * time.Second),
		},
		{
			Name:     "Another Day in Paradise",
			Duration: durationpb.New(123 * time.Second),
		},
	}

	var wg sync.WaitGroup
	wg.Add(3)
	for i := range audios {
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.CreateAudio(ctx, &grpcapi.CreateAudioRequest{Audio: &audios[i]})
			if err != nil {
				log.Fatalf("could not create audio: %v", err)
			}
			log.Printf("Audio added to playlist with id: %s", r.GetAudio().GetId())
		}(i)
	}
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Play(ctx, &grpcapi.PlayRequest{})
	if err != nil {
		log.Fatalf("could not play audio: %v", err)
	}

	time.Sleep(5 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Pause(ctx, &grpcapi.PauseRequest{})
	if err != nil {
		log.Fatalf("could not play audio: %v", err)
	}

	time.Sleep(2 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Play(ctx, &grpcapi.PlayRequest{})
	if err != nil {
		log.Fatalf("could not play audio: %v", err)
	}
}
