import React from 'react'
import {
  Avatar,
  Button,
  Icon,
  ListItemAvatar,
} from '@material-ui/core'
import { apiResetCheck } from './Api'
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
    apiResetCheck(
      this.props.id,
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
      {this.displayDialog({
        title: "Really reset check?",
        prompt: `Really reset the check named "${props.name}" in context "${props.context}"?`
      })}
      </>
  }
}
