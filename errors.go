package QuickHTTPServerLauncher

import "errors"

func enrichError(text string, err error) error {
	return errors.New(text + ": " + err.Error())
}
