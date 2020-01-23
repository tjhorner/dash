import React, { Component } from "react"
import axios from "axios"
import "./CommuteWidget.scss"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faSubway } from "@fortawesome/free-solid-svg-icons"

class CommuteWidget extends Component {
  constructor() {
    super()
    this.state = {
      time: null,
      refreshInterval: null
    }
  }

  async refreshTime() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/commute/time`)
    this.setState({ time: data.time })
  }

  componentDidMount() {
    const refreshInterval = setInterval(() => { this.refreshTime() }, 300000)
    this.setState({ refreshInterval })
    this.refreshTime()
  }

  componentWillUnmount() {
    clearInterval(this.state.refreshInterval)
  }
  
  render() {
    if(!this.state.time) return (<div></div>)

    return (
      <div className="widget widget-commute">
        <div className="widget-content">
          <div><FontAwesomeIcon icon={faSubway}/> <strong>Time to Work</strong></div>
          <div className="time">~{this.state.time}</div>
        </div>
      </div>
    )
  }
}

export default CommuteWidget