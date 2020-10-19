package QuickHTTPServerLauncher

import (
	"github.com/kaatinga/bufferedlogger"
	"os"
)

func NewConfig() (config Config) {

	// Setting up logger
	logger := bufferedlogger.InitLog(os.Stdout)
	config.logger = &logger
	return
}
