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
    if (name) {
      history.push(`/search/context:${context} name:${name}`)
    } else {
      history.push(`/search/context:${context}`)
    }
  }

  render() {
    const { name, context } = this.props
    let title
    if (name) {
      title = `Search for events with name "${name}" in context "${context}"`
    } else {
      title = `Search for events in context "${context}"`
    }
    return (
      <IconButton title={title} aria-label={title} onClick={this.handleClick}>
        <Search/>
      </IconButton>
    )
  }
}

