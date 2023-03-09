package player

type command uint8

const (
	Play command = iota
	Pause
	Next
	Prev
)

type commandMsg struct {
	command command
	err     chan error
}

type state uint8

const (
	noActiveAudio state = iota
	playing
	paused
	closed
)

type playerSignals struct {
	playCh  chan struct{}
	pauseCh chan struct{}
	endCh   chan struct{}
	closeCh chan struct{}
}
