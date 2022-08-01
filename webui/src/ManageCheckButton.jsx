import React from 'react'
import {
  IconButton,
} from '@material-ui/core'
import {
  Done
} from '@material-ui/icons'
import ManageCheck from './ManageCheck'

export default class ManageCheckButton extends React.Component {
  state = {
    open: false,
  }

  handleClick = () => {
    this.setState({ open: true })
  }

  handleClose = () => {
    this.setState({ open: false });
    this.props.refresh()
  }

  render() {
    const { eventName, eventContext, refresh } = this.props
    const title = "Manage check for this event"

    return (
      <>
        <IconButton title={title} aria-label={title} onClick={this.handleClick}>
          <Done/>
        </IconButton>
        <ManageCheck eventName={eventName} eventContext={eventContext} open={this.state.open} onClose={this.handleClose} refresh={refresh}/>
      </>
    )
  }
}
