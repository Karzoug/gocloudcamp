package server

import (
	"context"
	"errors"

	"github.com/Karzoug/gocloudcamp/internal/grpcapi"
	"github.com/Karzoug/gocloudcamp/internal/models"
	"github.com/Karzoug/gocloudcamp/internal/player"
	"github.com/Karzoug/gocloudcamp/internal/playlist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type server struct {
	grpcapi.PlayerServiceServer
	player *player.Player
}

func New(p *player.Player) *server {
	return &server{player: p}
}

func (s *server) Play(ctx context.Context, _ *grpcapi.PlayRequest) (*grpcapi.PlayResponse, error) {
	err := s.player.Play(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case errors.Is(err, player.ErrNoAudio):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.PlayResponse{}, nil
}
func (s *server) Pause(ctx context.Context, _ *grpcapi.PauseRequest) (*grpcapi.PauseResponse, error) {
	err := s.player.Pause(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.PauseResponse{}, err
}
func (s *server) Next(ctx context.Context, _ *grpcapi.NextRequest) (*grpcapi.NextResponse, error) {
	err := s.player.Next(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case errors.Is(err, player.ErrNoAudio):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.NextResponse{}, err
}
func (s *server) Prev(ctx context.Context, _ *grpcapi.PrevRequest) (*grpcapi.PrevResponse, error) {
	err := s.player.Prev(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case errors.Is(err, player.ErrNoAudio):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.PrevResponse{}, err
}
func (s *server) CreateAudio(ctx context.Context, req *grpcapi.CreateAudioRequest) (*grpcapi.CreateAudioResponse, error) {
	reqAudio := req.GetAudio()
	respAudio, err := s.player.Playlist.Add(ctx, models.Audio{
		Id:       reqAudio.GetId(),
		Name:     reqAudio.GetName(),
		Duration: reqAudio.GetDuration().AsDuration(),
	})
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case err == playlist.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.CreateAudioResponse{
		Audio: &grpcapi.Audio{
			Id:       respAudio.Id,
			Name:     respAudio.Name,
			Duration: durationpb.New(respAudio.Duration),
		},
	}, nil
}
func (s *server) ReadAudio(ctx context.Context, req *grpcapi.ReadAudioRequest) (*grpcapi.ReadAudioResponse, error) {
	reqAudioId := req.GetId()
	respAudio, err := s.player.Playlist.Get(ctx, reqAudioId)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case err == playlist.ErrNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.ReadAudioResponse{
		Audio: &grpcapi.Audio{
			Id:       respAudio.Id,
			Name:     respAudio.Name,
			Duration: durationpb.New(respAudio.Duration),
		},
	}, nil
}
func (s *server) UpdateAudio(ctx context.Context, req *grpcapi.UpdateAudioRequest) (*grpcapi.UpdateAudioResponse, error) {
	reqAudio := req.GetAudio()
	respAudio, err := s.player.Playlist.Update(ctx, models.Audio{
		Id:       reqAudio.GetId(),
		Name:     reqAudio.GetName(),
		Duration: reqAudio.GetDuration().AsDuration(),
	})
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case errors.Is(err, playlist.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, playlist.ErrCurrentAudio):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.UpdateAudioResponse{
		Audio: &grpcapi.Audio{
			Id:       respAudio.Id,
			Name:     respAudio.Name,
			Duration: durationpb.New(respAudio.Duration),
		},
	}, nil
}
func (s *server) DeleteAudio(ctx context.Context, req *grpcapi.DeleteAudioRequest) (*grpcapi.DeleteAudioResponse, error) {
	reqAudioId := req.GetId()
	err := s.player.Playlist.Delete(ctx, reqAudioId)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		case errors.Is(err, playlist.ErrCurrentAudio):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &grpcapi.DeleteAudioResponse{}, nil
}
func (s *server) ListAudio(ctx context.Context, _ *grpcapi.ListAudioRequest) (*grpcapi.ListAudioResponse, error) {
	slice, err := s.player.Playlist.List(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	resp := grpcapi.ListAudioResponse{
		Audio: []*grpcapi.Audio{},
	}
	for _, a := range slice {
		resp.Audio = append(resp.Audio, &grpcapi.Audio{
			Id:       a.Id,
			Name:     a.Name,
			Duration: durationpb.New(a.Duration),
		})
	}
	return &resp, nil
}
