import React from 'react'
import ReactDom from 'react-dom'
import {
  Chip,
  TableCell,
  TableRow,
  IconButton,
} from '@material-ui/core'
import {
  UnfoldLess,
  UnfoldMore,
  Share,
} from '@material-ui/icons'
import SearchButton from './SearchButton'
import { renderDuration, renderDate, renderCommandResult } from './DisplayHelpers'
import { apiGetEvent } from './Api'
import ManageCheckButton from './ManageCheckButton'
import { CopyToClipboard } from 'react-copy-to-clipboard'

export const commandString = (command) => {
  return command.map((c) => (c.includes(' ') ? ('"' + c + '"') : c)).join(' ')
}

const prettyJSON = require('json-align')

class EventsLinesOutput extends React.Component {
  formatJSON(json) {
    return prettyJSON(json, null, 2)
  }

  formatOutput(output) {
    let newOutput = ""
    if (!output) {
      return newOutput
    }
    output.split("\n").forEach( (l) => {
      if (l[0] !== "{") {
        newOutput += l + "\n"
        return
      }
      try {
        const json = JSON.parse(l)
        const newLine = this.formatJSON(json) + "\n"
        newOutput += newLine
      } catch(e) {
        newOutput += l + "\n"
      }
    })
    return newOutput
  }

  render() {
    const { active, hostname, command } = this.props

    if (!active) return null

    const commandResult = renderCommandResult(this.props)

    const output = this.formatOutput(this.props.outputLoaded)
    const preStyle = {
      margin: '5px',
      color: '#ffbf00',
      backgroundColor: '#301600',
      whiteSpace: 'pre-wrap',
      wordWrap: 'break-word'
    }
    return (
      <TableRow>
        <TableCell colSpan="7" style={{backgroundColor: 'black'}}>
          <pre style={preStyle}>
          @{hostname} $ {command && commandString(command) + ` ⮑ ${commandResult}\n`}
          {output}
          </pre>
        </TableCell>
      </TableRow>
    )
  }
}

class FoldButton extends React.Component {
  handleClick = () => {
    this.props.event.props.setSelectedId(null)
    this.props.event.loadOutput()
  }

  render() {
    const { active } = this.props

    if (active) {
      return (
        <IconButton title="Close" aria-label="Close" onClick={this.handleClick}>
          <UnfoldLess/>
        </IconButton>
      )
    } else {
      return (
        <IconButton title="Open" aria-label="Open" onClick={this.handleClick}>
          <UnfoldMore/>
        </IconButton>
      )
    }
  }
}

const ShareButton = ({ id }) => (
  <IconButton title='Copy link clipboard' aria-label='Copy link to clipboard'>
    <CopyToClipboard text={`${window.location.origin}/search/id:${id}`}>
      <Share/>
    </CopyToClipboard>
  </IconButton>
)


const Bar = "▁▂▃▄▅▆▇█"

const Load = ({ load }) => {
  return (
    <IconButton title={`${(100 * load).toFixed(2)}% Load`} aria-label={"Load"}>
      {Bar[Math.floor((Bar.length - 1) * load)]}
    </IconButton>
  )
}

export default class Event extends React.Component {
  state = {
    active: false,
    outputLoaded: "",
  }

  displayDuration() {
    return renderDuration(this.props.duration)
  }

  commandIsLong() {
		const command = this.props.command
    if (!command) {
      return false
    }
    const cmd = commandString(command)
    return cmd.length > 40
  }

  displayCommand() {
		const command = this.props.command
    if (!command) {
      return ''
    }
    const cmd = commandString(command)
    if (this.commandIsLong()) {
      let sliced = cmd.slice(0, 20)
      sliced += '…'
      return <tt title={cmd}>{sliced}</tt>
    } else {
      return <tt>{cmd}</tt>
    }
  }

  displayStyle(selected) {
    if (selected) {
      return {
        border: '3pt solid black'
      }
    }
  }

  scrollToRow() {
    if (this.props.selected){
      const selectedRow = ReactDom.findDOMNode(this.refs.row)
      window.scrollTo(0, selectedRow.offsetTop)
    }
  }

  componentDidUpdate() {
    this.scrollToRow()
  }

  loadOutput() {
    if (this.state.outputLoaded === "") {
      apiGetEvent(this.props.id, ({ data: { data } }) => {
        this.setState({ outputLoaded: data[0].output })
        this.setState((state) => ({ ...state, active: !state.active }))
      })
    } else {
      this.setState((state) => ({ ...state, active: !state.active }))
    }
  }

  render() {
		const {
      output, success, id, name, context, started, load, selected, refresh
    } = this.props

    return (
      <>
        <TableRow className={success ? 'success' : 'failure' } style={this.displayStyle(selected)} ref="row">
          <TableCell>
              <SearchButton context={context} name={name}/>
              <ManageCheckButton name={name} context={context} refresh={refresh}/>
              <ShareButton id={id}/>
              <Load load={load}/>
              <Chip label={name} color="primary"/>
          </TableCell>
          <TableCell>
            <SearchButton context={context} name={name}/>
            <Chip label={context} color="secondary"/>
          </TableCell>
          <TableCell>{renderDate(started)}</TableCell>
          <TableCell>{this.displayDuration()}</TableCell>
          <TableCell>
            <FoldButton active={this.state.active} event={this}/>
            {output}
          </TableCell>
        </TableRow>
      <EventsLinesOutput
        active={this.state.active}
        outputLoaded={this.state.outputLoaded}
        {...this.props}
      />
      </>
    )
  }
}
