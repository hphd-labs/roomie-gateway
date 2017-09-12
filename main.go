package main

import (
	"github.com/andrewburian/powermux"
	"net/http/httputil"
	"net/url"
	"os"
	"github.com/Sirupsen/logrus"
	"net/http"
)

const (
	ROUTE_AUTH = "auth"
)

func main() {

	// Logger
	logger := logrus.StandardLogger()

	// The auth gateway
	var authGateway *httputil.ReverseProxy
	{
		authUpstreamPath := os.Getenv("AUTH_UPSTREAM")
		authUpstream, err := url.Parse(authUpstreamPath)
		if err != nil {
			logger.WithField("auth_upstream", authUpstreamPath).Fatal("Invalid Auth Upstream")
		}
		authGateway = httputil.NewSingleHostReverseProxy(authUpstream)
	}


	// Setup the router
	mux := powermux.NewServeMux()
	mux.Route(ROUTE_AUTH).Any(authGateway)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "http"
	}
	err := http.ListenAndServe(":"+port, mux)
	logger.Fatal(err)

}