import React from 'react'
import { apiGetCheckByName, apiStoreCheck } from './Api'
import EditCheckDialog from './EditCheckDialog'

export default class ManageCheck extends React.Component {
  state = {
    check: null,
    newCheck: null
  }

  loadCheck() {
    const { eventName } = this.props
    apiGetCheckByName(eventName,
      ({ data: { data } }) => {
        this.setState({ newCheck: false, check: data[0] })
      },
      ({ response: { status } }) => {
        if (status === 404) {
          this.setState({ newCheck: true, check: { name: eventName } })
        }
      }
    )
  }

  handleCloseSave = check => {
    apiStoreCheck(check,
      () => {
        this.setState({ open: false, newCheck: null })
        this.props.onClose()
      }
    )
  }

  render() {
    const { open, onClose } = this.props
    const { newCheck, check } = this.state

    if (open) {
      if (newCheck === null) {
        this.loadCheck()
        return null
      } else if (newCheck === true) {
        return (
          <EditCheckDialog
            action="add"
            name={check.name}
            open={open}
            onClose={onClose}
            onCloseSave={this.handleCloseSave}
          />
        )
      } else if (newCheck === false) {
        return (
          <EditCheckDialog
            action="edit"
            open={open}
            onClose={onClose}
            onCloseSave={this.handleCloseSave}
            {...check}
          />
        )
      } else {
        return null
      }
    } else {
      return null
    }
  }
}

