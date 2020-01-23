package maps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type TravelMode string

const (
	TravelModeDriving   TravelMode = "driving"
	TravelModeWalking   TravelMode = "walking"
	TravelModeBicycling TravelMode = "bicycling"
	TravelModeTransit   TravelMode = "transit"
)

type TransitMode string

const (
	TransitModeBus    TransitMode = "bus"
	TransitModeSubway TransitMode = "subway"
	TransitModeTrain  TransitMode = "train"
	TransitModeTram   TransitMode = "tram"
	TransitModeRail   TransitMode = "rail"
)

type DistanceMatrixRequest struct {
	Key                      string      `url:"key"`
	Origins                  string      `url:"origins"`
	Destinations             string      `url:"destinations"`
	TravelMode               TravelMode  `url:"mode,omitempty`
	Language                 string      `url:"language,omitempty"`
	Region                   string      `url:"region,omitempty"`
	Avoid                    string      `url:"avoid,omitempty"`
	ArrivalTime              string      `url:"arrival_time,omitempty"`
	DepartureTime            string      `url:"departure_time,omitempty"`
	TrafficModel             string      `url:"traffic_model,omitempty"`
	TransitMode              TransitMode `url:"transit_mode,omitempty"`
	TransitRoutingPreference string      `url:"transit_routing_preference,omitempty"`
}

type DistanceMatrixValue struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type DistanceMatrixElement struct {
	Status   string              `json:"status"`
	Duration DistanceMatrixValue `json:"duration"`
	Distance DistanceMatrixValue `json:"distance"`
}

type DistanceMatrixRow struct {
	Elements []DistanceMatrixElement `json:"elements"`
}

type DistanceMatrixResponse struct {
	Status               string              `json:"status"`
	OriginAddresses      []string            `json:"origin_addresses"`
	DestinationAddresses []string            `json:"destination_addresses"`
	Rows                 []DistanceMatrixRow `json:"rows"`
}

// GetDistanceMatrix gets a distance matrix between one or more pairs of locations
func GetDistanceMatrix(params *DistanceMatrixRequest) (*DistanceMatrixResponse, error) {
	v, _ := query.Values(params)

	res, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?%s", v.Encode()))
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("http status: %s", res.Status)
	}

	defer res.Body.Close()
	j := json.NewDecoder(res.Body)
	var dm DistanceMatrixResponse
	err = j.Decode(&dm)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}
