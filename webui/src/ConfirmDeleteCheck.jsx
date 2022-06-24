import React from 'react'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
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
          <Dialog
            open={this.state.open}
            onClose={this.handleClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
          >
            <DialogTitle id="alert-dialog-title">{"Really delete check?"}</DialogTitle>
            <DialogContent>
              <DialogContentText id="alert-dialog-description">
                Really delete the check named "{name}"?
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
      </ListItemIcon>
    )
  }
}
