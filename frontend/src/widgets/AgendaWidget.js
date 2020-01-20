import React, { Component } from "react"
import "./AgendaWidget.scss"
import moment from "moment"
import axios from "axios"

class AgendaWidget extends Component {
  constructor() {
    super()
    this.state = {
      events: null
    }
  }

  async refreshEvents() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/agenda/events`)
    this.setState({ events: data })
  }

  componentDidMount() {
    this.refreshEvents()
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

    const events = this.state.events.map(event => (
      <div className="event" key={event.id}>
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