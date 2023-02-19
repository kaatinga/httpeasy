package httpeasy

import "github.com/julienschmidt/httprouter"

type Adapter func(httprouter.Handle) httprouter.Handle

// Adapt composes handlers chain end executes it.
func Adapt(next httprouter.Handle, adapters ...Adapter) httprouter.Handle {
	for _, adapter := range adapters {
		next = adapter(next)
	}
	return next
}
