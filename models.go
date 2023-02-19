package httpeasy

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

var timeOutDuration = 5 * time.Second

// SetUpHandlers type to announce handlers.
type SetUpHandlers func(r *httprouter.Router)

// Config - http service configuration compatible to settings package.
// https://github.com/kaatinga/settings
type Config struct {
	ProductionMode    bool          `env:"PROD"`
	SSL               SSL           `validate:"required_if=ProductionMode true"`
	Port              uint16        `env:"PORT" validate:"min=80,max=65535"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" default:"1m"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" default:"15s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" default:"1m"`
}

type SSL struct {
	Domain string `env:"DOMAIN" validate:"fqdn"`
	Email  string `env:"EMAIL" validate:"email"`
}

// newWebService creates http.Server structure with router inside.
func (config *Config) newWebService() http.Server {
	return http.Server{
		Addr:              net.JoinHostPort("", fmt.Sprintf("%d", config.Port)),
		Handler:           httprouter.New(),
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
	}
}

// Launch enables the configured web service with the handlers that
// announced in a function matched with SetUpHandlers type.
func (config *Config) Launch(handlers SetUpHandlers) error {
	webServer := config.newWebService()

	// enable handlers inside SetUpHandlers function
	router, ok := webServer.Handler.(*httprouter.Router)
	if !ok {
		return errors.New("webServer.Handler is not a *httprouter.Router")
	}
	handlers(router)

	// shutdown is a special channel to handle errors
	shutdown := make(chan error)

	switch config.ProductionMode {
	case true:
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,

			// Domain
			HostPolicy: autocert.HostWhitelist(config.SSL.Domain),

			// Folder to store certificates
			Cache: autocert.DirCache("certs"),
			Email: config.SSL.Email,
		}

		webServer.TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}

		// Config server to redirect
		go func() {
			_ = http.ListenAndServe(
				":http",
				certManager.HTTPHandler(

					// Redirect from http to https
					http.RedirectHandler(
						"https://"+config.SSL.Domain,
						http.StatusPermanentRedirect),
				),
			)
		}()

		// HTTPS server to handle the service
		go func() {
			funcErr := webServer.ListenAndServeTLS("", "")
			if funcErr != nil {
				shutdown <- funcErr
				close(shutdown)
			}
		}()
	default:
		go func() {
			funcErr := webServer.ListenAndServe()
			if funcErr != nil {
				shutdown <- funcErr
				close(shutdown)
			}
		}()
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var terminationSignal string
	select {
	case osSignal := <-interrupt:
		terminationSignal = osSignal.String()
	case err := <-shutdown:
		terminationSignal = err.Error()
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), timeOutDuration)
	defer cancelFunc()

	err := webServer.Shutdown(timeout)
	return fmt.Errorf("service stopped: %w, terminationSignal: %s", err, terminationSignal)
}
