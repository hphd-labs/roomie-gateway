package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ROUTE_AUTH = "/auth"
)

func main() {

	// Logger
	logger := logrus.StandardLogger()

	// Load config
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		logger.Fatal(err)
	}

	// The auth gateway
	authUpstream, err := url.Parse(conf.AuthUpstream)
	if err != nil {
		logger.WithField("AUTH_UPSTREAM", conf.AuthUpstream).Fatal("Invalid Auth Upstream")
	}
	authGateway := NewHostReverseProxy(authUpstream)

	// Setup the router
	mux := powermux.NewServeMux()
	mux.Route(ROUTE_AUTH).Any(authGateway).Route("*").Any(authGateway)

	// Setup the server
	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mux,
	}

	// Trap TERM and INT signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Signals kill the server
	go func(c <-chan os.Signal) {
		select {
		case sig := <-c:
			shutdownTime := time.Duration(conf.ShutdownTime) * time.Second
			shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), shutdownTime)
			logrus.WithField("signal", sig).Warn("Trapped signal")
			server.Shutdown(shutdownCtx)
			cancelFunc()
		}
	}(sigChan)

	// Run server
	err = server.ListenAndServe()

	// Clean exit on close
	if err == http.ErrServerClosed {
		return
	}
	logger.Fatal(err)

}
