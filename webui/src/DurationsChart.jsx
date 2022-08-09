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
        0: { title: 'Duration [m]', viewWindow: { min: 0 } },
        1: { textPosition: 'none' },
      },
      tooltip: {
        isHtml: true
      },
      trendlines: {
        0: { tooltip: false },
        1: { tooltip: false },
      },
      series: {
        0: { title: 'Duration', targetAxisIndex: 0, lineWidth: 0, pointSize: 5 },
        1: { title: 'Memory', targetAxisIndex: 1, lineWidth: 1, pointSize: 2, curveType: 'function', tooltip: false },
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
        id: 'Memory',
      },
    ]
  }

  prepareRows(data) {
    return data.map( (e) => {
      const { started, duration, memoryUsage } = e
      const startTime = Date.parse(started)
      return [
        new Date(startTime),
        duration / nano / 60,
        formatTooltip(e),
        e.success ? '#40f540' : '#f54040',
        memoryUsage,
      ]
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
      width="103%"
      chartEvents={chartEvents}
    />
  }
}
