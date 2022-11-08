import React from 'react'

import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from '@material-ui/core'

export default class Confirm extends React.Component {
  state = {
    open: false,
  }

  handleClickOpen = () => {
    this.setState({ open: true })
  }

  handleClose = () => {
    this.setState({ open: false })
  }

  handleCloseYes = () => {
    this.confirmAction()
    this.setState({ open: false })
  }

  displayDialog = ({ title, prompt }) => {
    return <Dialog
    open={this.state.open}
    onClose={this.handleClose}
    aria-labelledby="alert-dialog-title"
    aria-describedby="alert-dialog-description"
  >
      <DialogTitle id="alert-dialog-title">{title}</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          {prompt}
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
  }
}
