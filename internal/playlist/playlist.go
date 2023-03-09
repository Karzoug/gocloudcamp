package playlist

import (
	"context"

	"github.com/Karzoug/gocloudcamp/internal/models"
)

type Playlist interface {
	AudioRepository
	Current() *models.Audio
	CurrentToFront() *models.Audio
	CurrentToNext() *models.Audio
	CurrentToPrev() *models.Audio
	Front() *models.Audio
	Back() *models.Audio
}

type AudioRepository interface {
	Add(ctx context.Context, a models.Audio) (*models.Audio, error)
	Get(ctx context.Context, id string) (*models.Audio, error)
	Update(ctx context.Context, a models.Audio) (*models.Audio, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Audio, error)
	Close() error
}
