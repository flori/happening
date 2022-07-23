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
    history.push(`/search/${this.props.query}`)
  }

  render() {
    const title = "Search for events with this name"
    return (
      <IconButton title={title} aria-label={title} onClick={this.handleClick}>
        <Search/>
      </IconButton>
    )
  }
}

