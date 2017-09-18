package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	ROUTE_AUTH = "auth"
)

const (
	CONF_PORT          = "PORT"
	CONF_AUTH_UPSTREAM = "AUTH_UPSTREAM"
)

func main() {

	// Logger
	logger := logrus.StandardLogger()

	// The auth gateway
	var authGateway *httputil.ReverseProxy
	{
		authUpstreamPath := os.Getenv(CONF_AUTH_UPSTREAM)
		authUpstream, err := url.Parse(authUpstreamPath)
		if err != nil {
			logger.WithField(CONF_AUTH_UPSTREAM, authUpstreamPath).Fatal("Invalid Auth Upstream")
		}
		authGateway = httputil.NewSingleHostReverseProxy(authUpstream)
	}

	// Setup the router
	mux := powermux.NewServeMux()
	mux.Route(ROUTE_AUTH).Any(authGateway).Route("*").Any(authGateway)

	// Run the server
	port := os.Getenv(CONF_PORT)
	if port == "" {
		port = "http"
	}
	err := http.ListenAndServe(":"+port, mux)
	logger.Fatal(err)

}
