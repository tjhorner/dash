package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tjhorner/dash/util"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// SpotifyAPI is the API that handles the Spotify panel
type SpotifyAPI struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthToken    string
	RefreshToken string

	spotifyClient *spotify.Client
}

// Configure implements API.Configure
func (api *SpotifyAPI) Configure() {
	api.ClientID = util.GetEnv("DASH_SPOTIFY_CLIENT_ID", "")
	api.ClientSecret = util.GetEnv("DASH_SPOTIFY_CLIENT_SECRET", "")
	api.RedirectURL = util.GetEnv("DASH_SPOTIFY_REDIRECT_URL", "")
	api.AuthToken = util.GetEnv("DASH_SPOTIFY_AUTH_TOKEN", "")
	api.RefreshToken = util.GetEnv("DASH_SPOTIFY_REFRESH_TOKEN", "")

	token := &oauth2.Token{
		Expiry:       time.Now(), // so we always refresh the token
		TokenType:    "Bearer",
		AccessToken:  api.AuthToken,
		RefreshToken: api.RefreshToken,
	}
	auth := spotify.NewAuthenticator(api.RedirectURL, spotify.ScopeUserReadPlaybackState)
	auth.SetAuthInfo(api.ClientID, api.ClientSecret)
	client := auth.NewClient(token)
	api.spotifyClient = &client
}

// Prefix implements API.Prefix
func (api *SpotifyAPI) Prefix() string {
	return "/api/spotify"
}

// Route implements API.Route
func (api *SpotifyAPI) Route(router *mux.Router) {
	router.HandleFunc("/playing", api.getPlaying).Methods("GET")
}

// GET /playing
func (api *SpotifyAPI) getPlaying(w http.ResponseWriter, r *http.Request) {
	playing, err := api.spotifyClient.PlayerState()
	if err != nil {
		respondError(err, "ERR_SPOTIFY_RESPONSE", w, r)
		return
	}

	respondJSON(*playing, w, r)
}
