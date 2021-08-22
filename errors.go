package httpeasy

import cerr "github.com/kaatinga/const-errs"

const (
	ErrRouterTypeIsIncorrect = cerr.Error("Incorrect Router Type")
	ErrLoggerIsNotEnabled    = cerr.Error("Incorrect Config Logger")
)
