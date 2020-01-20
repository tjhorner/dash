import React, { Component } from "react"
import axios from "axios"
import "./SpotifyWidget.scss"
import defaultCover from "../default-cover-art.png"

class SpotifyWidget extends Component {
  constructor() {
    super()
    this.state = {
      playing: null,
      refreshInterval: null
    }
  }

  async refreshPlaying() {
    const { data } = await axios.get(`${process.env.REACT_APP_API_BASE_URL || ""}/api/spotify/playing`)
    this.setState({ playing: data })
  }

  componentDidMount() {
    const refreshInterval = setInterval(() => { this.refreshPlaying() }, 10000)
    this.setState({ refreshInterval })
    this.refreshPlaying()
  }

  componentWillUnmount() {
    clearInterval(this.state.refreshInterval)
  }
  
  render() {
    if(!this.state.playing) return (<div></div>)

    if(!this.state.playing.is_playing) {
      return (
        <div className="widget widget-spotify">
          <div className="widget-content">
            <div className="album-art" style={{ backgroundImage: `url("${defaultCover}")`  }}/>
            <div className="track-details">
              <div><em>Nothing is playing on Spotify</em></div>
            </div>
          </div>
        </div>
      )
    }

    return (
      <div className="widget widget-spotify">
        <div className="widget-content">
          <div className="album-art" style={{ backgroundImage: `url("${this.state.playing.Item.album.images[0].url}")` }}/>
          <div className="track-details">
            <div className="device muted">Playing on <strong>{this.state.playing.device.name}</strong></div>
            <div><strong>{this.state.playing.Item.name}</strong></div>
            <div>{this.state.playing.Item.artists.length > 1 ? "Various Artists" : this.state.playing.Item.artists[0].name}</div>
          </div>
        </div>
      </div>
    )
  }
}

export default SpotifyWidget