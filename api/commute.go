package api

import (
	"net/http"
	"time"

	"github.com/tjhorner/dash/api/citymapper"

	"github.com/gorilla/mux"
	"github.com/tjhorner/dash/util"
)

// CommuteAPI is the API that handles the Commute panel
type CommuteAPI struct {
	CitymapperKey string
	FromCoords    string
	ToCoords      string

	latestTime citymapper.TravelTime
	timeFrom   time.Time
}

// Configure implements API.Configure
func (api *CommuteAPI) Configure() {
	api.CitymapperKey = util.GetEnv("DASH_COMMUTE_CITYMAPPER_KEY", "")
	api.FromCoords = util.GetEnv("DASH_COMMUTE_FROM_COORDS", "")
	api.ToCoords = util.GetEnv("DASH_COMMUTE_TO_COORDS", "")
}

// Prefix implements API.Prefix
func (api *CommuteAPI) Prefix() string {
	return "/api/commute"
}

// Route implements API.Route
func (api *CommuteAPI) Route(router *mux.Router) {
	router.HandleFunc("/time", api.getTime).Methods("GET")
}

// GET /time
func (api *CommuteAPI) getTime(w http.ResponseWriter, r *http.Request) {
	// simple caching
	if time.Now().Sub(api.timeFrom) <= time.Hour {
		respondJSON(api.latestTime, w, r)
		return
	}

	tt, err := citymapper.GetTravelTime(api.FromCoords, api.ToCoords, api.CitymapperKey)
	if err != nil {
		respondError(err, "ERR_CITYMAPPER_RESPONSE", w, r)
		return
	}

	api.latestTime = *tt
	api.timeFrom = time.Now()
	respondJSON(*tt, w, r)
}
