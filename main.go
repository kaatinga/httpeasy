package QuickHTTPServerLauncher

import (
	"github.com/kaatinga/bufferedlogger"
	"os"
)

func NewConfig() (config Config) {

	// Setting up the config logger
	logger := bufferedlogger.InitLog(os.Stdout)
	config.Logger = &logger
	return
}
