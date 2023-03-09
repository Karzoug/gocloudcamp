package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Karzoug/gocloudcamp/internal/models"
	"github.com/Karzoug/gocloudcamp/internal/playlist/memory"
)

type filePlaylistConfig interface {
	StoreFile() string
	Restore() bool
}

type FilePlaylist struct {
	memory.MemPlaylist
	cfg filePlaylistConfig
}

func New(cfg filePlaylistConfig) (*FilePlaylist, error) {
	fp := &FilePlaylist{
		MemPlaylist: *memory.New(),
		cfg:         cfg,
	}

	if cfg.Restore() {
		if err := fp.restore(); err != nil {
			return nil, fmt.Errorf("restore playlist from file error: %w", err)
		}
	}

	return fp, nil
}

func (fp *FilePlaylist) Close() error {
	return fp.saveData()
}

func (fp *FilePlaylist) saveData() error {
	log.Printf("Save filelist to file: %s", fp.cfg.StoreFile())

	auds, err := fp.List(context.TODO())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fp.cfg.StoreFile(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open store file error: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(auds)
	return err
}

func (fp *FilePlaylist) restore() error {
	log.Printf("Restore filelist from file: %s", fp.cfg.StoreFile())

	file, err := os.OpenFile(fp.cfg.StoreFile(), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	auds := make([]models.Audio, 0)
	if err := json.NewDecoder(file).Decode(&auds); err != nil && err != io.EOF {
		return err
	}
	if err := fp.MemPlaylist.SetAll(auds); err != nil {
		return err
	}
	return nil
}
