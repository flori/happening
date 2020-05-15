import React from 'react'
import EventsChart from './EventsChart'
import DurationsChart from './DurationsChart'
import Event from './Event'
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from '@material-ui/core'
import { Subscribe } from 'unstated'

class EventsListComponent extends React.Component {
  state = {
    selectedId: null,
  }

  setSelectedId = id => { this.setState({ selectedId: id }) }

  get selectedId() {
    return this.state.selectedId
  }

  computeChartHeight(events) {
    let height = 0

    const eventCount = {}
    events.forEach( (e) => { eventCount[e.name] = true } )

    height = 64 * Object.keys(eventCount).length

    if (height > 512) {
      height = 512
    }

    return height
  }

  render() {
    const { events, loaded, refresh } = this.props

    const height = this.computeChartHeight(events)

    return (
      <Paper style={{margin: '0.1em'}}>
        <EventsChart loaded={loaded} height={height} key={height} data={events} setSelectedId={this.setSelectedId} />
        <DurationsChart loaded={loaded} data={events} setSelectedId={this.setSelectedId} />
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Started</TableCell>
              <TableCell>Duration</TableCell>
              <TableCell>Output</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
          {events.map(
            e => <Event key={e.id} selected={this.selectedId === e.id} {...e} setSelectedId={this.setSelectedId} refresh={refresh}/>
          )}
          </TableBody>
        </Table>
      </Paper>
    )
  }
}

export default class EventsList extends React.Component {
  get eventsContainer() {
    return this.props.eventsContainer
  }

  get query() {
    return this.props.match.params.query
  }

  get seconds() {
    return this.eventsContainer.seconds
  }

  componentDidMount() {
    let seconds = this.seconds
    const params = new URLSearchParams(this.props.location.search)
    const sParam = params.get('s')
    if (sParam) {
      seconds = parseInt(sParam, 10)
    }
    this.eventsContainer.load(this.query, seconds)
  }

  componentDidUpdate() {
    this.eventsContainer.load(this.query, this.seconds)
  }

  render() {
    return (
      <Subscribe to={[this.eventsContainer]}>
        {ec => <EventsListComponent events={ec.events} loaded={ec.loaded} refresh={this.props.refresh}/>}
      </Subscribe>
    )
  }
}
