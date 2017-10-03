package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"net/url"
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

	err = http.ListenAndServe(":"+conf.Port, mux)
	if err == http.ErrServerClosed {
		return
	}
	logger.Fatal(err)

}
