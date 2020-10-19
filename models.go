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

	"github.com/kaatinga/bufferedlogger"
	"golang.org/x/crypto/acme/autocert"
)

type SetUpHandlers func(r *httprouter.Router, db *sql.DB)

type Config struct {
	Email      string
	DB         *sql.DB
	LaunchMode string
	Port       string
	Domain     string
	Logger     *bufferedlogger.Logger
}

func (config *Config) Launch(handlers SetUpHandlers) error {
	var err error
	config.Logger.Title.Info().Str("port", config.Port).Msg("Launching the service on the")

	// Create a new router
	var router = httprouter.New()

	// Configuration validation
	if config.DB == nil {
		return errors.New("the DB connection is not ok #1")
	}

	handlers(router, config.DB)
	config.Logger.SubMsg.Info().Msg("Handlers have been announced")

	var webServer http.Server

	switch config.LaunchMode {
	case "prod":
		config.Logger.SubMsg.Info().Msg("Production Mode is enabled")
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,

			// Domain
			HostPolicy: autocert.HostWhitelist(config.Domain),

			// Folder for storing certificates
			Cache: autocert.DirCache("certs"),
			Email: config.Email,
		}

		webServer = http.Server{
			Addr:              net.JoinHostPort("", config.Port),
			Handler:           router,
			ReadTimeout:       1 * time.Minute,
			ReadHeaderTimeout: 15 * time.Second,
			WriteTimeout:      1 * time.Minute,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go func() {
			go func() {
				err := http.ListenAndServe(
					":http",
					certManager.HTTPHandler(

						// Redirect from http to https
						http.RedirectHandler(
							strings.Join([]string{"https://", config.Domain}, ""),
							http.StatusPermanentRedirect),
					),
				)
				if err != nil {
					config.Logger.SubMsg.Err(err).Msg("Service stopped")
				}
			}()

			err := webServer.ListenAndServeTLS("", "")
			if err != nil {
				config.Logger.SubMsg.Err(err).Msg("Service stopped")
			}
		}()
	case "dev":
		config.Logger.SubMsg.Warn().Msg("Development Mode is enabled")

		webServer = http.Server{
			Addr:              net.JoinHostPort("", config.Port),
			Handler:           router,
			ReadTimeout:       1 * time.Minute,
			ReadHeaderTimeout: 15 * time.Second,
			WriteTimeout:      1 * time.Minute,
		}

		go func() {
			err := webServer.ListenAndServe()
			if err != nil {
				config.Logger.SubMsg.Err(err).Msg("Service stopped")
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
		config.Logger.SubMsg.Err(err).Msg("Service stopped")
	}

	return nil
}
