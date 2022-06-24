import React from 'react'

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
}
