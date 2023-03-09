package memory

import (
	"container/list"
	"context"
	"log"
	"sync"

	"github.com/Karzoug/gocloudcamp/internal/models"
	"github.com/Karzoug/gocloudcamp/internal/playlist"
	"github.com/rs/xid"
)

type MemPlaylist struct {
	list    *list.List
	current *list.Element
	mtx     sync.RWMutex
}

func New() *MemPlaylist {
	p := MemPlaylist{
		list: list.New(),
		mtx:  sync.RWMutex{},
	}
	p.current = p.list.Front()
	return &p
}

func (p *MemPlaylist) Current() *models.Audio {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.current == nil {
		return nil
	}
	audio, _ := p.current.Value.(models.Audio)
	return &audio
}

func (p *MemPlaylist) CurrentToFront() *models.Audio {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.current = p.list.Front()
	if p.current == nil {
		return nil
	}
	audio, _ := p.current.Value.(models.Audio)
	return &audio
}

func (p *MemPlaylist) CurrentToNext() *models.Audio {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.current == nil {
		return nil
	}
	p.current = p.current.Next()
	if p.current == nil {
		return nil
	}
	audio, _ := p.current.Value.(models.Audio)
	return &audio
}

func (p *MemPlaylist) CurrentToPrev() *models.Audio {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.current == nil {
		return nil
	}
	p.current = p.current.Prev()
	if p.current == nil {
		return nil
	}
	audio, _ := p.current.Value.(models.Audio)
	return &audio
}

func (p *MemPlaylist) Front() *models.Audio {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.list.Front() == nil {
		return nil
	}
	audio, _ := p.list.Front().Value.(models.Audio)
	return &audio
}

func (p *MemPlaylist) Back() *models.Audio {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if p.list.Back() == nil {
		return nil
	}
	audio, _ := p.list.Back().Value.(models.Audio)
	return &audio
}

func (p *MemPlaylist) Add(_ context.Context, a models.Audio) (*models.Audio, error) {
	log.Printf("Add new audio: %s", a.Name)

	p.mtx.Lock()
	defer p.mtx.Unlock()

	a.Id = xid.New().String()
	p.list.PushBack(a)

	return &a, nil
}

func (p *MemPlaylist) Get(_ context.Context, id string) (*models.Audio, error) {
	log.Printf("Get audio with id: %s", id)

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	for e := p.list.Front(); e != nil; e = e.Next() {
		if v := e.Value.(models.Audio); v.Id == id {
			return &v, nil
		}
	}
	return nil, playlist.ErrNotFound
}

func (p *MemPlaylist) Update(_ context.Context, a models.Audio) (*models.Audio, error) {
	log.Printf("Update audio with id: %s", a.Id)

	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.current != nil && p.current.Value.(models.Audio).Id == a.Id {
		return nil, playlist.ErrCurrentAudio
	}

	for e := p.list.Front(); e != nil; e = e.Next() {
		if v := e.Value.(models.Audio); v.Id == a.Id {
			e.Value = a
			v = e.Value.(models.Audio)
			return &v, nil
		}
	}
	return nil, playlist.ErrNotFound
}

func (p *MemPlaylist) Delete(_ context.Context, id string) error {
	log.Printf("Delete audio with id: %s", id)

	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.current != nil && p.current.Value.(models.Audio).Id == id {
		return playlist.ErrCurrentAudio
	}

	for e := p.list.Front(); e != nil; e = e.Next() {
		if v := e.Value.(models.Audio); v.Id == id {
			p.list.Remove(e)
			return nil
		}
	}
	return nil
}

func (p *MemPlaylist) List(_ context.Context) ([]models.Audio, error) {
	log.Println("Get audio list")

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	slice := make([]models.Audio, 0, p.list.Len())
	for e := p.list.Front(); e != nil; e = e.Next() {
		slice = append(slice, e.Value.(models.Audio))
	}
	return slice, nil
}

func (p *MemPlaylist) Close() error {
	return nil
}
