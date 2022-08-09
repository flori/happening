import React from 'react'
import { Chart } from 'react-google-charts'
import LinearProgress from '@material-ui/core/LinearProgress';
import { formatTooltip, micro } from './DisplayHelpers'

export default class EventsChart extends React.Component {
  state = {
    options: {
      timeline: {
        colorByRowLabel: true,
        groupByRowLabel: true,
      },
      hAxis: {
        format: 'M-d H:mm',
      },
      tooltip: {
        isHtml: true
      }
    },
    columns: [
      {
        type: 'string',
        id: 'Name',
      },
      {
        type: 'string',
        id: 'Id',
      },
      {
        type: 'string',
        role: 'tooltip',
        p: { html: true },
      },
      {
        type: 'date',
        id: 'Start',
      },
      {
        type: 'date',
        id: 'End',
      },
    ],
  }

  prepareRows(data) {
    return data.map( (e) => {
      const { name, started, duration } = e
      const startTime = Date.parse(started)
      const endTime = startTime + duration / micro
      const tooltip = formatTooltip(e)
      return [ name, '', tooltip, new Date(startTime), new Date(endTime) ]
    })
  }

  setupChartEvents(data, setSelectedId) {
    return [
      {
        name: 'select',
        callback({ chartWrapper }) {
          const event = data[chartWrapper.getChart().getSelection()[0].row]
          if (event) {
            setSelectedId(event.id)
          }
        },
      },
    ]
  }

  render() {
    const { data, height, loaded, setSelectedId } = this.props
    if (loaded < 1) return <LinearProgress size={50} />
    if (data.length < 2) return null

    const rows = this.prepareRows(data)
    const chartEvents = this.setupChartEvents(data, setSelectedId)

    return <Chart
      chartType="Timeline"
      rows={rows}
      columns={this.state.columns}
      options={this.state.options}
      graph_id="Timeline"
      width="99%"
      height={`${height}px`}
      chartEvents={chartEvents}
    />
  }
}
