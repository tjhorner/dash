package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/option"

	"golang.org/x/oauth2"

	"google.golang.org/api/calendar/v3"

	"github.com/gorilla/mux"
	"github.com/tjhorner/dash/util"
)

// AgendaAPI is the API that handles the Agenda panel
type AgendaAPI struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleAuthToken    string
	GoogleRefreshToken string

	apiClient *calendar.Service
}

// Configure implements API.Configure
func (api *AgendaAPI) Configure() {
	api.GoogleClientID = util.GetEnv("DASH_AGENDA_GOOGLE_CLIENT_ID", "")
	api.GoogleClientSecret = util.GetEnv("DASH_AGENDA_GOOGLE_CLIENT_SECRET", "")
	api.GoogleRedirectURL = util.GetEnv("DASH_AGENDA_GOOGLE_REDIRECT_URL", "")
	api.GoogleAuthToken = util.GetEnv("DASH_AGENDA_GOOGLE_AUTH_TOKEN", "")
	api.GoogleRefreshToken = util.GetEnv("DASH_AGENDA_GOOGLE_REFRESH_TOKEN", "")

	config := &oauth2.Config{
		ClientID:     api.GoogleClientID,
		ClientSecret: api.GoogleClientSecret,
		RedirectURL:  api.GoogleRedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v4/token",
		},
	}
	token := &oauth2.Token{
		Expiry:       time.Now(),
		TokenType:    "Bearer",
		AccessToken:  api.GoogleAuthToken,
		RefreshToken: api.GoogleRefreshToken,
	}
	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		fmt.Printf("warn: unable to set up google calendar service (%v)\n", err)
	}

	api.apiClient = calendarService
}

// Prefix implements API.Prefix
func (api *AgendaAPI) Prefix() string {
	return "/api/agenda"
}

// Route implements API.Route
func (api *AgendaAPI) Route(router *mux.Router) {
	router.HandleFunc("/events", api.getEvents).Methods("GET")
}

// GET /events
func (api *AgendaAPI) getEvents(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	year, month, day := t.Date()

	bod := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	eod := time.Date(year, month, day, 23, 59, 59, 0, t.Location())

	events, err := api.apiClient.Events.List("primary").TimeMin(bod.Format(time.RFC3339)).TimeMax(eod.Format(time.RFC3339)).Do()
	if err != nil {
		respondError(err, "ERR_CALENDAR_RESPONSE", w, r)
		return
	}

	respondJSON(events.Items, w, r)
}
