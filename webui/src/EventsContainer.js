import { Container } from 'unstated'
import { apiGetEvents } from './Api'

export default class EventsContainer extends Container {
  state = {
    events: [],
    loaded: 0,
    seconds: 3600,
    query: null,
    updated: 0
  }

  get events() {
    return this.state.events
  }

  get loaded() {
    return this.state.loaded
  }

  get seconds() {
    return this.state.seconds
  }

  set seconds(newSeconds) {
    this.load(this.state.query, newSeconds)
  }

  load(query, seconds) {
    this.setState({ loaded: 0 })
    apiGetEvents({ query, seconds }, ({ data: { data } }) => {
      this.setState(
        state => {
          return {
            query, seconds, loaded: state.loaded + 1, events: data || [],
          }
      })
    })
  }
}
