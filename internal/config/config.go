package config

import (
	"errors"
	"flag"
)

type Config struct {
	storeFile string
	restore   bool
}

const (
	defaultStoreFile = "/tmp/gocloud_player.json"
	defaultRestore   = true
)

// New creates Config with default values.
func New() *Config {
	return &Config{
		storeFile: defaultStoreFile,
		restore:   defaultRestore,
	}
}

func (c Config) StoreFile() string {
	return c.storeFile
}

func (c Config) Restore() bool {
	return c.restore
}

func (с Config) IsStoreInMemory() bool {
	return с.storeFile == ""
}

// Load loads flags config values to Config.
func (c *Config) Load() error {
	if flag.Parsed() {
		return errors.New("flags have already been parsed")
	}

	flag.StringVar(&c.storeFile, "f", defaultStoreFile, "filename to save/load playlist")
	flag.BoolVar(&c.restore, "r", defaultRestore, "whether to load saved data at startup")
	flag.Parse()

	return nil
}
