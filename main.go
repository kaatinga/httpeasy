package httpeasy

import (
	"github.com/kaatinga/bufferedlogger"
	"os"
)

func NewConfig() (config Config) {

	// Setting up the config Logger
	logger := bufferedlogger.InitLog(os.Stdout)
	config.Logger = &logger
	return
}
