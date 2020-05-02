import React from 'react'
import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  TextField,
} from '@material-ui/core'
import { renderDuration, renderDate, parseDuration, nano } from './DisplayHelpers'
import { CheckStateIcon } from './checkState'

export default class EditCheckDialog extends React.Component {
  constructor(props) {
    super(props)
    let { id, name, period, disabled, allowed_failures } = this.props
    if (!period) {
      period = 3900 * nano
    }
    this.state = { id, name, period, disabled, allowed_failures }
  }

  handleNameChange = e => {
    const name = e.target.value
    if (name) {
      this.setState({ name })
    }
  }

  handlePeriodChange = e => {
    const period = parseDuration(e.target.value)
    if (period) {
      this.setState({ period })
    }
  }

  handleAllowedFailuersChange = e => {
    let allowed_failures = parseInt(e.target.value, 10)
    if (isNaN(allowed_failures) || allowed_failures < 0) {
      allowed_failures = 0
    }
    this.setState({ allowed_failures })
  }

  handleDisabledChange = e => {
    let disabled = e.target.checked
    this.setState({ disabled })
  }

  render() {
    const { action, last_ping_at, open, failures, success, healthy } = this.props
    const { id, name, period, disabled, allowed_failures } = this.state
    const actionName = action.replace(/\b\w/g, l => l.toUpperCase())

    return (
      <Dialog
        open={open}
        onClose={this.props.onClose}
        aria-labelledby="form-dialog-title"
      >
        <DialogTitle id="form-dialog-title">
          <span style={{paddingRight: '0.25em'}}>{actionName} Check "{name}"</span>
          <CheckStateIcon action={action} disabled={disabled} healthy={healthy} success={success}/>
        </DialogTitle>
        <DialogContent>
          <form noValidate>
            <TextField
              autoFocus={name == null}
              margin="dense"
              id="name"
              label="Name"
              type="text"
              value={name}
              disabled={action === "edit"}
              onChange={this.handleNameChange}
              fullWidth
            />
            {last_ping_at != null &&
              <TextField
                margin="dense"
                id="last-ping-at"
                label="Last Ping At"
                type="text"
                value={renderDate(last_ping_at)}
                disabled
                fullWidth
              />}
            <TextField
              autoFocus={name != null}
              margin="dense"
              id="failures"
              label="Failures"
              type="number"
              value={failures}
              disabled
              fullWidth
            />
            <TextField
              autoFocus={name != null}
              margin="dense"
              id="allowed-failures"
              name="allowed_failures"
              label="Allowed Failures"
              type="number"
              defaultValue={this.state.allowed_failures}
              onChange={this.handleAllowedFailuersChange}
              inputProps={{min: 0}}
              fullWidth
            />
            <TextField
              autoFocus={name != null}
              margin="dense"
              id="period"
              name="period"
              label="Period"
              defaultValue={renderDuration(this.state.period)}
              onChange={this.handlePeriodChange}
              fullWidth
            />
            <FormControlLabel
              label="Disabled"
              control={
                <Checkbox
                  margin="dense"
                  id="disabled"
                  checked={disabled}
                  onChange={this.handleDisabledChange}
                />
              }
            />
          </form>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Cancel
          </Button>
          <Button onClick={() => this.props.onCloseSave({ id, name, disabled, period, allowed_failures })} color="primary">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    )
  }
}
