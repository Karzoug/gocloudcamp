package player

import "errors"

var (
	ErrNoAudio      = errors.New("no audio to play")
	ErrPlayerClosed = errors.New("player closed")
)
