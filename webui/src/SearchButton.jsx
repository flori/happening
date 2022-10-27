import React from 'react'
import {
  IconButton,
} from '@material-ui/core'
import {
  Search
} from '@material-ui/icons'
import { history } from './history'

export default class SearchButton extends React.Component {
  handleClick = () => {
    const { context, name } = this.props
    history.push(`/search/context:${context} name:${name}`)
  }

  render() {
    const { name, context } = this.props
    const title = `Search for events with name ${name} in ${context}`
    return (
      <IconButton title={title} aria-label={title} onClick={this.handleClick}>
        <Search/>
      </IconButton>
    )
  }
}

