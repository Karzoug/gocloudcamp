package player

import (
	"context"
	"fmt"
	"log"

	"github.com/Karzoug/gocloudcamp/internal/models"
	"github.com/Karzoug/gocloudcamp/internal/playlist"
	"github.com/ivahaev/timer"
)

var mockPlayFnc = func(a models.Audio, signals playerSignals) error {
	log.Printf("New audio loaded: name '%s', duration '%s'", a.Name, a.Duration)
	var (
		err error
		t   *timer.Timer
	)

	go func() {

		select {
		case <-signals.closeCh:
			log.Print("Audio closed")
			return
		case <-signals.playCh:
			log.Print("Audio started")
			t = timer.NewTimer(a.Duration)
			t.Start()
			defer t.Stop()
		}

		for {
			select {
			case <-signals.closeCh:
				log.Print("Audio closed")
				return
			case <-t.C:
				log.Print("Audio ended")
				signals.endCh <- struct{}{}
				return
			case <-signals.pauseCh:
				log.Print("Audio paused")
				t.Pause()
			case <-signals.playCh:
				log.Print("Audio started")
				t.Start()
			}
		}
	}()

	return err
}

type Player struct {
	Playlist playlist.Playlist
	playFnc  func(models.Audio, playerSignals) error

	commandsCh chan commandMsg
	state      state
	signals    playerSignals

	closePlayerCh chan struct{}
}

func New(pl playlist.Playlist) *Player {
	p := Player{
		Playlist:   pl,
		commandsCh: make(chan commandMsg, 10),
		signals: playerSignals{
			playCh:  make(chan struct{}),
			pauseCh: make(chan struct{}),
			closeCh: make(chan struct{}),
			endCh:   make(chan struct{}),
		},
		closePlayerCh: make(chan struct{}),
		playFnc:       mockPlayFnc,
	}
	if a := p.Playlist.Current(); a != nil {
		if err := p.handleCurrentElement(); err != nil {
			p.state = noActiveAudio
		} else {
			p.state = paused
		}

	}

	go p.loop()

	return &p
}

func (p *Player) Close() error {
	close(p.closePlayerCh)
	return p.Playlist.Close()
}

// Play начинает воспроизведение
func (p *Player) Play(ctx context.Context) error {
	return p.addCommand(ctx, Play)
}

// Pause приостанавливает воспроизведение
func (p Player) Pause(ctx context.Context) error {
	return p.addCommand(ctx, Pause)
}

// Next позволяет воспроизвести след песню
func (p *Player) Next(ctx context.Context) error {
	return p.addCommand(ctx, Next)
}

// Prev позволяет воспроизвести предыдущую песню
func (p *Player) Prev(ctx context.Context) error {
	return p.addCommand(ctx, Prev)
}

func (p Player) addCommand(ctx context.Context, c command) error {
	errCh := make(chan error)
	msg := commandMsg{
		command: c,
		err:     errCh,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.closePlayerCh:
		return ErrPlayerClosed
	case p.commandsCh <- msg:
	}

	select {
	case <-ctx.Done():
		go func() {
			<-errCh
		}()
		return ctx.Err()
	case <-p.closePlayerCh:
		return ErrPlayerClosed
	case err := <-errCh:
		return err
	}
}

func (p *Player) loop() {
	for {
		select {
		case <-p.closePlayerCh:
			p.state = closed
			close(p.signals.closeCh)
			return
		case c := <-p.commandsCh:
			switch c.command {
			case Play:
				p.play(c.err)
			case Pause:
				p.pause(c.err)
			case Next:
				p.next(c.err)
			case Prev:
				p.prev(c.err)
			}
		case <-p.signals.endCh:
			p.state = noActiveAudio
			errCh := make(chan error)
			go func() {
				<-errCh
			}()
			p.next(errCh)
		}
	}
}

func (p *Player) play(errCh chan error) {
	switch p.state {
	case paused:
	case noActiveAudio:
		if p.Playlist.Current() == nil {
			p.Playlist.CurrentToFront()
		}
		if err := p.handleCurrentElement(); err != nil {
			errCh <- err
			return
		}
	default:
		errCh <- nil
		return
	}

	p.signals.playCh <- struct{}{}
	p.state = playing
	errCh <- nil
}

func (p *Player) pause(errCh chan error) {
	if p.state != playing {
		errCh <- nil
		return
	}
	p.signals.pauseCh <- struct{}{}
	p.state = paused
	errCh <- nil
}

func (p *Player) next(errCh chan error) {
	switch p.state {
	case playing, paused:
		p.signals.closeCh <- struct{}{}
		p.state = noActiveAudio
	case noActiveAudio:
	default:
		errCh <- nil
		return
	}

	p.Playlist.CurrentToNext()
	if err := p.handleCurrentElement(); err != nil {
		errCh <- err
		return
	}
	p.state = paused
	p.signals.playCh <- struct{}{}
	p.state = playing
	errCh <- nil
}

func (p *Player) prev(errCh chan error) {
	switch p.state {
	case playing, paused:
		p.signals.closeCh <- struct{}{}
		p.state = noActiveAudio
	case noActiveAudio:
	default:
		errCh <- nil
		return
	}

	p.Playlist.CurrentToPrev()
	if err := p.handleCurrentElement(); err != nil {
		errCh <- err
		return
	}
	p.state = paused
	p.signals.playCh <- struct{}{}
	p.state = playing
	errCh <- nil
}

func (p *Player) handleCurrentElement() error {
	if p.Playlist.Current() == nil {
		return ErrNoAudio
	}
	if err := p.playFnc(*p.Playlist.Current(), p.signals); err != nil {
		return fmt.Errorf("handle audio problem: %w", err)
	}
	return nil
}
