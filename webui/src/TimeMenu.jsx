import React from 'react'
import { IconButton, Menu, MenuItem } from '@material-ui/core'
import { Watch } from '@material-ui/icons'
import { renderDuration, nano } from './DisplayHelpers'
import { Subscribe } from 'unstated'

const TimeChoice = ({ menu, ec, chosenSeconds, text }) => {
  const seconds = ec.seconds
  return (
    <MenuItem
      selected={seconds === chosenSeconds}
      onClick={menu.handleChoice(ec, chosenSeconds)}
    >{text}
    </MenuItem>
  )
}

export default class TimeMenu extends React.Component {
  state = { anchorEl: null }

  get eventsContainer() {
    return this.props.eventsContainer
  }

  handleClick = event => {
    event.persist()
    this.setState({ anchorEl: event.target })
  }

  handleChoice = (ec, seconds) => {
    return event => {
      if (seconds) {
        ec.seconds = seconds
      }
      this.setState({ anchorEl: null })
    }
  }

  render() {
    const { anchorEl } = this.state
    const title = "Duration"
    return (
      <Subscribe to={[this.eventsContainer]}>
        {ec =>
          <div title={renderDuration(ec.seconds * nano)}>
            <IconButton
              title={title} aria-label={title}
              aria-owns={anchorEl ? 'time-menu' : null}
              aria-haspopup="true"
              onClick={this.handleClick}
            >
              <Watch/>
            </IconButton>
            <Menu
              id="time-menu"
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={this.handleChoice(ec, null)}
            >
              <TimeChoice menu={this} ec={ec} chosenSeconds={         1 * 3600} text='1h'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={         2 * 3600} text='2h'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={         6 * 3600} text='6h'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={        12 * 3600} text='12h'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={    1 * 24 * 3600} text='1d'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={    2 * 24 * 3600} text='2d'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={    3 * 24 * 3600} text='3d'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={1 * 7 * 24 * 3600} text='1w'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={2 * 7 * 24 * 3600} text='2w'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={3 * 7 * 24 * 3600} text='3w'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={4 * 7 * 24 * 3600} text='4w'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={8 * 7 * 24 * 3600} text='8w'/>
              <TimeChoice menu={this} ec={ec} chosenSeconds={16 * 7 * 24 * 3600} text='16w'/>
            </Menu>
          </div>
        }
      </Subscribe>
    )
  }
}
