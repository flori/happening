import React from 'react'
import {
  IconButton,
  ListItemIcon,
} from '@material-ui/core'
import DeleteIcon from '@material-ui/icons/Delete'
import { apiDeleteCheck } from './Api'
import Confirm from './Confirm'

export default class ConfirmDeleteCheck extends Confirm {
  confirmAction = () => {
    apiDeleteCheck(
      { id: this.props.id },
      this.props.refresh
    )
    this.setState({ open: false })
  }

  render() {
    const { name } = this.props
    return (
      <ListItemIcon>
        <>
          <IconButton aria-label="Delete" onClick={this.handleClickOpen}>
             <DeleteIcon/>
          </IconButton>
          {this.displayDialog({
            title: "Really delete check?",
            prompt: `Really delete the check named "${name}"?`,
          })}
        </>
      </ListItemIcon>
    )
  }
}
