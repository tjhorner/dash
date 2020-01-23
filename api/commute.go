package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tjhorner/dash/api/maps"
	"github.com/tjhorner/dash/util"
)

// CommuteAPI is the API that handles the Commute panel
type CommuteAPI struct {
	MapsKey     string
	FromAddress string
	ToAddress   string

	latestTime *timeResponse
	timeFrom   time.Time
}

type timeResponse struct {
	Time string `json:"time"`
}

// Configure implements API.Configure
func (api *CommuteAPI) Configure() {
	api.MapsKey = util.GetEnv("DASH_COMMUTE_GMAPS_KEY", "")
	api.FromAddress = util.GetEnv("DASH_COMMUTE_FROM_ADDRESS", "")
	api.ToAddress = util.GetEnv("DASH_COMMUTE_TO_ADDRESS", "")
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
	// simple caching (1min)
	if time.Now().Sub(api.timeFrom) <= time.Minute {
		respondJSON(api.latestTime, w, r)
		return
	}

	dmr := &maps.DistanceMatrixRequest{
		Key:          api.MapsKey,
		Origins:      api.FromAddress,
		Destinations: api.ToAddress,
		TravelMode:   maps.TravelModeTransit,
		TransitMode:  maps.TransitModeSubway,
	}
	dm, err := maps.GetDistanceMatrix(dmr)
	if err != nil {
		respondError(err, "ERR_MAPS_RESPONSE", w, r)
		return
	}

	resp := timeResponse{
		Time: dm.Rows[0].Elements[0].Duration.Text,
	}

	api.latestTime = &resp
	api.timeFrom = time.Now()
	respondJSON(resp, w, r)
}
