package playlist

import "errors"

var (
	ErrNotFound     = errors.New("audio not found")
	ErrCurrentAudio = errors.New("invalid argument: this is the current audio")
)
