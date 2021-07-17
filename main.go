package httpeasy

import (
	"github.com/kaatinga/prettylogger"
	"github.com/rs/zerolog"
)

// NewConfig creates new config model that later can be processed
// by settings package or updated manually.
func NewConfig(logLevel zerolog.Level, colour bool) (config *Config) {

	config = new(Config)

	// Setting up the config Logger
	config.Logger = prettylogger.InitLogger(logLevel, colour, true)
	return
}
