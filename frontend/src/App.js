import React from "react"
import "./App.css"
import "normalize.css"
import ClimateWidget from "./widgets/ClimateWidget"
import SpotifyWidget from "./widgets/SpotifyWidget"
import CommuteWidget from "./widgets/CommuteWidget"
import ClockWidget from "./widgets/ClockWidget"
import AgendaWidget from "./widgets/AgendaWidget"

function App() {
  return (
    <div className="widgets">
      <div className="column">
        <SpotifyWidget/>
        <ClimateWidget/>
      </div>
      <div className="column">
        <ClockWidget/>
        <CommuteWidget/>
      </div>
      <div className="column">
        <AgendaWidget/>
      </div>
    </div>
  );
}

export default App;
