package bolt // import "github.com/justanotherorganization/justanotherbotkit/users/bolt"

import (
	"errors"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// Config is a database config.
type Config struct {
	// File is the bolt db path/name.
	File string
	// Options are the bolt db options.
	Options *bolt.Options
}

// Validate a configuration, applying defaults where possible.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("cfg cannot be nil")
	}

	if strings.TrimSpace(c.File) == "" {
		return errors.New("cfg.File cannot be empty")
	}

	return nil
}
