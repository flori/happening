import React from 'react'
import {
  Avatar,
  Icon,
  ListItemAvatar,
} from '@material-ui/core'

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

export const CheckStateAvatar = (props) => 
  <ListItemAvatar>
    <Avatar title={props.title} aria-label={props.title}>
      <CheckStateIcon {...props}/>
    </Avatar>
  </ListItemAvatar>
