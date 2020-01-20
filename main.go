package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/tjhorner/dash/api"

	"github.com/gobuffalo/packr/v2"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	router := mux.NewRouter()

	routeAPIs(
		router,
		&api.ClimateAPI{},
		&api.SpotifyAPI{},
		&api.CommuteAPI{},
		&api.AgendaAPI{},
	)

	frontendBox := packr.New("frontend", "./frontend/build")

	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		index, err := frontendBox.Find("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Write(index)
	}

	serveSpa := func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !frontendBox.Has(path) {
			serveIndex(w, r)
			return
		}

		http.FileServer(frontendBox).ServeHTTP(w, r)
	}

	// http.FileServer will redirect index.html to / and will end up in a redirect loop.
	// So we need to do this to fix that redirect loop.
	router.HandleFunc("/", serveIndex)
	router.PathPrefix("/").HandlerFunc(serveSpa)

	fmt.Printf("dash is now listening at %s\n", envListenAddr)

	log.Fatal(http.ListenAndServe(envListenAddr, handlers.CORS()(router)))
}
