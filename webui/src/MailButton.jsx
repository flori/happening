import React from 'react'
import PropTypes from 'prop-types'
import {
  IconButton,
} from '@material-ui/core'
import {
  Mail,
} from '@material-ui/icons'
import { apiPatchMailEvent } from './Api'

class MailButton extends React.Component {
  handleClick = () => {
    const { id } = this.props
    apiPatchMailEvent(id)
  }

  render() {
    return <IconButton title='Send notification mail' aria-label='Send notification mail' onClick={this.handleClick}>
      <Mail/>
    </IconButton>
  }
}

export default MailButton
