package main

import (
	"github.com/gorilla/mux"
)

// API implements an API interface
type API interface {
	Configure()
	Prefix() string
	Route(*mux.Router)
}

func routeAPIs(r *mux.Router, a ...API) {
	for _, api := range a {
		routeAPI(api, r)
	}
}

func routeAPI(a API, r *mux.Router) {
	a.Configure()
	sr := r.PathPrefix(a.Prefix()).Subrouter()
	a.Route(sr)
}
