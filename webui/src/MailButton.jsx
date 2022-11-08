import React from 'react'
import PropTypes from 'prop-types'
import {
  IconButton,
} from '@material-ui/core'
import {
  Mail,
} from '@material-ui/icons'
import { apiPatchMailEvent } from './Api'
import Confirm from './Confirm'

class MailButton extends Confirm {
  confirmAction = () => {
    apiPatchMailEvent(this.props.id)
  }

  render() {
    const { context, name } = this.props
    return <>
      <IconButton title='Send notification mail' aria-label='Send notification mail' onClick={this.handleClickOpen}>
        <Mail/>
      </IconButton>
      {this.displayDialog({
        title: "Send mail fo event",
        prompt: `Really send a mail for event "${name}" on "${context}?"`,
      })}
      </>
  }
}

export default MailButton
