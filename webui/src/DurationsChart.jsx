import React from 'react'
import { Chart } from 'react-google-charts'
import LinearProgress from '@material-ui/core/LinearProgress';
import { formatTooltip, nano } from './DisplayHelpers'

export default class DurationsChart extends React.Component {
  state = {
    options: {
      legend: 'none',
      hAxis: {
        format: 'M-d H:mm'
      },
      vAxes: {
        0: { title: 'Duration [m]', viewWindow: { min: 0, } },
        1: { title: 'Load [0-100%]', viewWindow: { min: 0, } },
      },
      tooltip: {
        isHtml: true
      },
      trendlines: {
        0: { tooltip: false },
        1: null,
      },
      series: {
        0: { title: 'Duration', targetAxisIndex: 0, lineWidth: 0, pointSize: 5 },
        1: { title: 'Load', targetAxisIndex: 1, lineWidth: 1, pointSize: 0, curveType: 'function', tooltip: false },
      },
    },
    columns: [
      {
        type: 'date',
        id: 'Start',
      },
      {
        type: 'number',
        id: 'Duration',
      },
      {
        type: 'string',
        role: 'tooltip',
        p: { html: true },
      },
      {
        type: 'string',
        role: 'style',
      },
      {
        type: 'number',
        id: 'Load',
      },
    ]
  }

  prepareRows(data) {
    return data.map( (e) => {
      const { started, duration, load } = e
      const startTime = Date.parse(started)
      return [
        new Date(startTime),
        duration / nano / 60,
        formatTooltip(e),
        e.success ? '#00d700' : '#f71000',
        100 * load,
      ]
    })

  }

  setupChartEvents(data, setSelectedId) {
    return [
      {
        eventName: 'select',
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
    const { data, setSelectedId, loaded } = this.props
    if (loaded < 1) return <LinearProgress size={50} />
    if (data.length < 2) return null

    const rows = this.prepareRows(data)
    const chartEvents = this.setupChartEvents(data, setSelectedId)

    return <Chart
      chartType="LineChart"
      rows={rows}
      columns={this.state.columns}
      options={this.state.options}
      graph_id="Durations"
      width="100%"
      chartEvents={chartEvents}
    />
  }
}
