import React from 'react'
import {
  Avatar,
  Button,
  Icon,
  ListItemAvatar,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from '@material-ui/core'
import { apiPatchCheck } from './Api'
import Confirm from './Confirm'

export const checkState = ({ disabled, healthy, success }) => {
  let state = "check"
  let stateColor = "primary"
  let stateDescription = "healthy"
  if (disabled) {
    state = "clear"
    stateColor = "disabled"
    stateDescription = "disabled"
  } else {
    if (!healthy) {
      stateColor = "error"
      state = "sentiment_very_dissatisfied"
      stateDescription = "failed"
      if (success) {
        state = "timer"
        stateDescription = "timeout"
      }
    }
  }
  return { state, stateColor, stateDescription }
}

export const CheckStateIcon = ({ action, disabled, healthy, success }) => {
  if (action === "add") return null
  const { state, stateColor } = checkState({ disabled, healthy, success })
  return <Icon color={stateColor}>{state}</Icon>
}

export class CheckStateAvatar extends Confirm {
  confirmAction = () => {
    const patchedCheck = {
      ...this.props,
      ...{ failures: 0, success: true, healthy: true  }
    }
    apiPatchCheck(
      this.props.id,
      patchedCheck,
      this.props.refresh
    )
  }

  render() {
    const props = this.props
    return <>
      <ListItemAvatar onClick={this.handleClickOpen}>
        <Avatar title={props.title} aria-label={props.title} onClick={this.handleClick}>
          <CheckStateIcon {...props}/>
        </Avatar>
      </ListItemAvatar>
      <Dialog
        open={this.state.open}
        onClose={this.handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
          <DialogTitle id="alert-dialog-title">{"Really reset check?"}</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              Really reset the check named "{props.name}"?
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleClose} color="primary" autoFocus>
              No
            </Button>
            <Button onClick={this.handleCloseYes} color="primary">
              Yes
            </Button>
          </DialogActions>
        </Dialog>
      </>
  }
}
