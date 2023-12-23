package httpeasy

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
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

	terminate  chan struct{}
	terminated chan struct{}
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

func (config *Config) Init() {
	config.terminate = make(chan struct{})
	config.terminated = make(chan struct{})
}

func (config *Config) Terminate() {
	config.terminate <- struct{}{}
	<-config.terminated
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
			MinVersion:     tls.VersionTLS13,
		}

		// Config server to redirect
		go func() {
			_ = http.ListenAndServe( //nolint:gosec
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
			}
		}()
	default:
		go func() {
			funcErr := webServer.ListenAndServe()
			if funcErr != nil {
				shutdown <- funcErr
			}
		}()
	}

	<-config.terminate

	timeout, cancelFunc := context.WithTimeout(context.Background(), timeOutDuration)
	defer cancelFunc()

	var outputError error

	if err := webServer.Shutdown(timeout); err != nil {
		err := webServer.Close()
		if err != nil {
			return err
		}
		outputError = fmt.Errorf("unable to terminate the web service: %w", err)
	}
	if outputError != nil {
		return fmt.Errorf("%w: web service terminated: %w", outputError, <-shutdown)
	} else {
		outputError = fmt.Errorf("web service terminated: %w", <-shutdown)
	}

	return outputError
}
