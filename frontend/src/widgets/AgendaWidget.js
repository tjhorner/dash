import React, { Component } from "react"
import "./AgendaWidget.scss"
import moment from "moment"
import axios from "axios"

class AgendaWidget extends Component {
  constructor() {
    super()
    this.state = {
      events: null,
      refreshInterval: null
    }
  }

  async refreshEvents() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/agenda/events`)
    this.setState({ events: data })
  }

  componentDidMount() {
    const refreshInterval = setInterval(() => { this.refreshEvents() }, 60000)
    this.setState({ refreshInterval })
    this.refreshEvents()
  }

  componentWillUnmount() {
    clearInterval(this.state.refreshInterval)
  }
  
  render() {
    if(!this.state.events) return (<div></div>)

    if(this.state.events.length === 0) {
      return (
        <div className="widget widget-agenda">
          <div className="widget-content">
            <div className="title">Today's Agenda</div>
            <div className="no-events">
              <em>Nothing planned today :)</em>
            </div>
          </div>
        </div>
      )
    }

    const events = this.state.events.filter(event => {
      const isActuallyToday = moment(event.start.dateTime).isBetween(moment().startOf("day"), moment().endOf("day"))
      const declined = event.attendees.find(at => at.self === true && at.responseStatus === "declined")
      return isActuallyToday && !declined
    }).map(event => (
      <div className={`event ${moment().isAfter(moment(event.end.dateTime)) ? "past" : ""}`} key={event.id}>
        <div className="event-time">{moment(event.start.dateTime).format("h:mm A")} - {moment(event.end.dateTime).format("h:mm A")}</div>
        <div className="event-title">{event.summary}</div>
      </div>
    ))

    return (
      <div className="widget widget-agenda">
        <div className="widget-content">
          <div className="title">Today's Agenda</div>
          <div className="events">
            {events}
          </div>
        </div>
      </div>
    )
  }
}

export default AgendaWidget