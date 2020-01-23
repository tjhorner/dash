package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tjhorner/dash/api/thermostatd"

	forecast "github.com/mlbright/darksky/v2"

	"github.com/gorilla/mux"
	"github.com/tjhorner/dash/util"
)

// ClimateAPI is the API that handles the Climate panel
type ClimateAPI struct {
	WeatherLocationLat string
	WeatherLocationLng string
	DarkSkyKey         string
	ThermostatdKey     string
	ThermostatdHost    string

	latestForecast forecast.Forecast
	forecastFrom   time.Time
}

// Configure implements API.Configure
func (api *ClimateAPI) Configure() {
	api.WeatherLocationLat = util.GetEnv("DASH_CLIMATE_LOCATION_LAT", "")
	api.WeatherLocationLng = util.GetEnv("DASH_CLIMATE_LOCATION_LNG", "")
	api.DarkSkyKey = util.GetEnv("DASH_CLIMATE_DARK_SKY_KEY", "")
	api.ThermostatdKey = util.GetEnv("DASH_CLIMATE_THERMOSTATD_KEY", "")
	api.ThermostatdHost = util.GetEnv("DASH_CLIMATE_THERMOSTATD_HOST", "")
}

// Prefix implements API.Prefix
func (api *ClimateAPI) Prefix() string {
	return "/api/climate"
}

// Route implements API.Route
func (api *ClimateAPI) Route(router *mux.Router) {
	router.HandleFunc("/weather", api.getWeather).Methods("GET")
	router.HandleFunc("/thermostat", api.getThermostat).Methods("GET")
}

// GET /weather
func (api *ClimateAPI) getWeather(w http.ResponseWriter, r *http.Request) {
	// simple caching
	if time.Now().Sub(api.forecastFrom) <= time.Hour {
		respondJSON(api.latestForecast.Currently, w, r)
		return
	}

	currentForecast, err := forecast.Get(api.DarkSkyKey, api.WeatherLocationLat, api.WeatherLocationLng, "now", forecast.US, forecast.English)
	if err != nil || currentForecast.Code >= 400 {
		if err == nil {
			err = fmt.Errorf("response code from Dark Sky was not successful (%d)", currentForecast.Code)
			fmt.Printf("%+v\n", currentForecast)
		}
		respondError(err, "ERR_FORECAST_RESPONSE", w, r)
		return
	}

	api.forecastFrom = time.Now()
	api.latestForecast = *currentForecast
	respondJSON(currentForecast.Currently, w, r)
}

// GET /thermostat
func (api *ClimateAPI) getThermostat(w http.ResponseWriter, r *http.Request) {
	state, err := thermostatd.GetState(api.ThermostatdHost, api.ThermostatdKey)
	if err != nil {
		respondError(err, "ERR_THERMOSTATD_RESPONSE", w, r)
		return
	}

	respondJSON(*state, w, r)
}
