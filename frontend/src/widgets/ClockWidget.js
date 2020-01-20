import React, { Component } from "react"
import moment from "moment"
import "./ClockWidget.scss"

class ClockWidget extends Component {
  constructor() {
    super()
    this.state = {
      time: null,
      refreshInterval: null
    }
  }

  refreshTime() {
    this.setState({
      time: moment()
    })
  }

  componentDidMount() {
    const refreshInterval = setInterval(() => { this.refreshTime() }, 1000)
    this.setState({ refreshInterval })
    this.refreshTime()
  }

  componentWillUnmount() {
    clearInterval(this.state.refreshInterval)
  }
  
  render() {
    if(!this.state.time) return (<div></div>)

    return (
      <div className="widget widget-clock">
        <div className="widget-content">
          <div className="time">{this.state.time.format("h:mm:ss A")}</div>
          <div>{this.state.time.format("dddd, MMMM D, YYYY")}</div>
        </div>
      </div>
    )
  }
}

export default ClockWidget