package QuickHTTPServerLauncher

import (
	"context"
	"crypto/tls"
	"database/sql"
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
	Name       string
	Email      string
	DB         *sql.DB
	LaunchMode string
	Port       string
	Domain     string
	logger     *bufferedlogger.Logger
}

func (config *Config) Launch(handlers SetUpHandlers) {
	var err error

	config.logger.Title.Info().Msg(strings.Join([]string{"=====", config.Name, " service is starting  ====="}, ""))

	// создаём роутер
	var router = httprouter.New()

	handlers(router, config.DB)
	config.logger.SubMsg.Info().Msg("Handlers have been announced")

	var webServer http.Server

	switch config.LaunchMode {
	case "prod":
		config.logger.SubMsg.Info().Msg("Production Mode is enabled")
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
					config.logger.SubMsg.Err(err).Msg("Service stopped")
				}
			}()

			err := webServer.ListenAndServeTLS("", "")
			if err != nil {
				config.logger.SubMsg.Err(err).Msg("Service stopped")
			}
		}()
	case "dev":
		config.logger.SubMsg.Warn().Msg("Development Mode is enabled")

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
				config.logger.SubMsg.Err(err).Msg("Service stopped")
			}
		}()
	default:
		config.logger.SubMsg.Error().Msg("Incorrect Launch Mode")
	}

	config.logger.SubMsg.Info().Msg("The service has been launched!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err = webServer.Shutdown(timeout)
	if err != nil {
		config.logger.SubMsg.Err(err).Msg("Service stopped")
	}
}
