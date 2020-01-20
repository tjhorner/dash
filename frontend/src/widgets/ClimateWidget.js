import React, { Component } from "react"
import axios from "axios"
import "./ClimateWidget.scss"

class ClimateWidget extends Component {
  constructor() {
    super()
    this.state = {
      weather: null,
      thermostat: null,
      refreshInterval: null
    }
  }

  async refreshWeather() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/climate/weather`)
    this.setState({ weather: data })
  }

  async refreshThermostat() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/climate/thermostat`)
    this.setState({ thermostat: data })
  }

  refreshAll() {
    this.refreshWeather()
    this.refreshThermostat()
  }

  componentDidMount() {
    const refreshInterval = setInterval(() => { this.refreshAll() }, 10000)
    this.setState({ refreshInterval })
    this.refreshAll()
  }

  componentWillUnmount() {
    clearInterval(this.state.refreshInterval)
  }
  
  render() {
    if(!this.state.weather || !this.state.thermostat) return (<div></div>)

    return (
      <div className="widget widget-climate">
        <div className="widget-content">
          <div className="weather-icon">
            <i className={`wi wi-forecast-io-${this.state.weather.icon}`}></i>
          </div>
          <div className="weather-details">
            <div><strong>Brooklyn, NY</strong></div>
            <div>{Math.round(this.state.weather.temperature)}&deg;F <span className="muted">({Math.round(this.state.weather.apparentTemperature)}&deg;F)</span> and {this.state.weather.summary}</div>
            <div><strong>Thermostat:</strong> {this.state.thermostat.target_temperature}&deg;F ({this.state.thermostat.current_mode})</div>
          </div>
        </div>
      </div>
    )
  }
}

export default ClimateWidget