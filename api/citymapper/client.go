package citymapper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// TravelTime is a response to a travel time request
type TravelTime struct {
	TravelTimeMinutes int `json:"travel_time_minutes"`
}

type travelTimeRequestParameters struct {
	Start string `url:"startcoord"`
	End   string `url:"endcoord"`
	Key   string `url:"key"`
}

// GetTravelTime gets the travel time between 2 coordinates
func GetTravelTime(start string, end string, key string) (*TravelTime, error) {
	params := travelTimeRequestParameters{
		Start: start,
		End:   end,
		Key:   key,
	}
	v, _ := query.Values(params)

	res, err := http.Get(fmt.Sprintf("https://developer.citymapper.com/api/1/traveltime/?%s", v.Encode()))
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("http status: %s", res.Status)
	}

	defer res.Body.Close()
	j := json.NewDecoder(res.Body)
	var tt TravelTime
	err = j.Decode(&tt)
	if err != nil {
		return nil, err
	}

	return &tt, nil
}
