package httpeasy

import (
	"github.com/kaatinga/bufferedlogger"
	"os"
)

// NewConfig creates new config model that later can be processed
// by settings package or updated manually.
func NewConfig() (config *Config) {

	config = new(Config)

	// Setting up the config Logger
	logger := bufferedlogger.InitLog(os.Stdout)
	config.Logger = &logger
	return
}
