package QuickHTTPServerLauncher

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kaatinga/assets"
	"github.com/kaatinga/bufferedlogger"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/go-playground/validator.v9"
)

// Function type for announcing handlers.
type SetUpHandlers func(r *httprouter.Router, db *sql.DB)

type Config struct {
	email      string
	DB         *sql.DB
	launchMode string
	port       string
	domain     string
	Logger     *bufferedlogger.Logger
	hasDB      bool
}

// logStopped logs the reason why the web service was stopped.
func (config *Config) logStopped(reason error) {
	config.Logger.SubMsg.Err(reason).Msg("Service stopped")
}

func (config *Config) SetEmail(email string) {
	config.email = email
}

func (config *Config) SetLaunchMode(mode string) {
	config.launchMode = mode
}

func (config *Config) SetPort(port string) {
	config.port = port
}

func (config *Config) SetDomain(domain string) {
	config.domain = domain
}

func (config *Config) SetDBMode() {
	config.hasDB = true
}

// check validates the web service configuration.
func (config *Config) check() error {

	v := validator.New()

	err := v.Var(config.domain, "fqdn")
	if err != nil {
		return err
	}

	err = v.Var(config.email, "email")
	if err != nil {
		return err
	}

	if !(config.launchMode == "prod" || config.launchMode == "dev") {
		return errors.New("incorrect launch mode is set")
	}

	port := assets.CheckUint16(config.port)
	if !port.Ok {
		return errors.New("incorrect port number")
	}

	if port.Parameter < 1001 || port.Parameter > 9999 {
		return errors.New("incorrect port range")
	}

	if config.hasDB && config.DB == nil {
		return errors.New("the DB connection is nil")
	}

	return nil
}

// Launch enables the configured web service with the handlers that
// announced in a function matched with SetUpHandlers type.
func (config *Config) Launch(handlers SetUpHandlers) error {
	var err error

	// Configuration Validation
	err = config.check()
	if err != nil {
		return err
	}

	// Launching
	config.Logger.Title.Info().Str("port", config.port).Msg("Launching the service on the")

	// Create a new router
	var router = httprouter.New()

	handlers(router, config.DB)
	config.Logger.SubMsg.Info().Msg("Handlers have been announced")

	var webServer http.Server

	switch config.launchMode {
	case "prod":
		config.Logger.SubMsg.Info().Msg("Production Mode is enabled")
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,

			// domain
			HostPolicy: autocert.HostWhitelist(config.domain),

			// Folder for storing certificates
			Cache: autocert.DirCache("certs"),
			Email: config.email,
		}

		webServer = http.Server{
			Addr:              net.JoinHostPort("", config.port),
			Handler:           router,
			ReadTimeout:       1 * time.Minute,
			ReadHeaderTimeout: 15 * time.Second,
			WriteTimeout:      1 * time.Minute,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		// HTTP server to redirect
		go func() {
			err := http.ListenAndServe(
				":http",
				certManager.HTTPHandler(

					// Redirect from http to https
					http.RedirectHandler(
						strings.Join([]string{"https://", config.domain}, ""),
						http.StatusPermanentRedirect),
				),
			)
			if err != nil {
				config.logStopped(err)
			}
		}()

		// HTTPS server to handle the service
		go func() {
			err := webServer.ListenAndServeTLS("", "")
			if err != nil {
				config.logStopped(err)
			}
		}()
	case "dev":
		config.Logger.SubMsg.Warn().Msg("Development Mode is enabled")

		webServer = http.Server{
			Addr:              net.JoinHostPort("", config.port),
			Handler:           router,
			ReadTimeout:       1 * time.Minute,
			ReadHeaderTimeout: 15 * time.Second,
			WriteTimeout:      1 * time.Minute,
		}

		go func() {
			err := webServer.ListenAndServe()
			if err != nil {
				config.logStopped(err)
			}
		}()
	default:
		config.Logger.SubMsg.Error().Msg("Incorrect Launch Mode")
	}

	config.Logger.SubMsg.Info().Msg("The service has been launched!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err = webServer.Shutdown(timeout)
	if err != nil {
		config.logStopped(err)
	}

	return nil
}
