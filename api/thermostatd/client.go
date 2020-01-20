package thermostatd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// State is a thermostat state
type State struct {
	// PoweredOn is true if the thermostat is on
	PoweredOn bool `json:"powered_on"`
	// CurrentMode is the current mode the thermostat is in
	CurrentMode Mode `json:"current_mode"`
	// FanSpeed is the current fan speed
	FanSpeed FanSpeed `json:"fan_speed"`
	// TargetTemperature is the target temperature in Fahrenheit
	TargetTemperature int `json:"target_temperature"`
	// CurrentTemperature is the current temperature in Fahrenheit as represented by the thermometer (currently not in use)
	CurrentTemperature int `json:"current_temperature"`
}

// Mode is a mode that the thermostat can be in
type Mode string

const (
	// ModeCool is the cooling mode
	ModeCool Mode = "COOL"
	// ModeDry is the "dry" mode
	ModeDry Mode = "DRY"
	// ModeHeat is the heating mode
	ModeHeat Mode = "HEAT"
	// ModeFan is the fan-only mode
	ModeFan Mode = "FAN"
)

// FanSpeed describes a speed at which the fan can run at
type FanSpeed string

const (
	// FanSpeedAuto will make the A/C determine what speed the fan should run at automatically
	FanSpeedAuto FanSpeed = "AUTO"
	// FanSpeedQuiet is the lowest fan speed (1/4)
	FanSpeedQuiet FanSpeed = "QUIET"
	// FanSpeedLow is the second lowest fan speed (2/4)
	FanSpeedLow FanSpeed = "LOW"
	// FanSpeedMedium is the second highest fan speed (3/4)
	FanSpeedMedium FanSpeed = "MEDIUM"
	// FanSpeedHigh is the highest fan speed (4/4)
	FanSpeedHigh FanSpeed = "HIGH"
)

// GetState gets the state of a thermostatd instance
func GetState(host string, key string) (*State, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/state", host), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("http status: %s", res.Status)
	}

	defer res.Body.Close()
	j := json.NewDecoder(res.Body)
	var state State
	err = j.Decode(&state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}
