import React from 'react'
import { Link } from 'react-router-dom'
import {
  AppBar,
  IconButton,
  TextField,
  Toolbar,
  Tooltip,
  Typography,
  withStyles,
} from '@material-ui/core'
import {
  Search,
  Refresh,
  DoneAll,
  ExitToApp,
  Clear
} from '@material-ui/icons'
import { getAuth, clearAuth, apiURLPrefix } from './Api'
import TimeMenu from './TimeMenu'
import { history } from './history'

const styles = {
  root: {
    display: 'flex',
    flexDirection: 'row',
  },
  title: {
    flex: '1 1 0',
    marginRight: '1em',
    flexDirection: 'row',
    paddingRight: '0.5em',
  },
  form: {
    flex: '10 1 0',
    flexDirection: 'row',
  },
  formTextField: {
    width: '100%',
    marginBottom: '0.5em',
  },
}

const RefreshButton = ({refresh}) => {
  const title = 'Refresh'
  return (
    <IconButton title={title} aria-label={title} onClick={refresh}>
      <Refresh/>
    </IconButton>
  )
}

const ChecksButton = () => {
  const title = 'Manage checks'
  return (
    <IconButton title={title} aria-label={title} component={Link} to="/checks">
      <DoneAll/>
    </IconButton>
  )
}

class LoginLogout extends React.Component {
  handleLogout() {
    clearAuth()
    window.location.href = "/"
  }

  handleLogin() {
    clearAuth()
    window.location.href = "/login"
  }

  render() {
    if (getAuth()) {
      const title = "Logout"
      return (
        <IconButton title={title} aria-label={title} style={{marginLeft: '0.5em'}} onClick={this.handleLogout}>
          <ExitToApp/>
        </IconButton>
      )
    } else {
      const title = "Login"
      return (
        <IconButton title={title} aria-label={title} style={{marginLeft: '0.5em'}} onClick={this.handleLogin}>
          <ExitToApp/>
        </IconButton>
      )
    }
  }
}

const SearchPrefix = ({refresh}) => (
  <div style={{position: 'absolute', left: -40, top: 10, width: 20, height: 20}}>
    <Tooltip id="tooltip-icon" title="Possible key:value filters are: id, name, output, hostname, command, success, ec">
      <IconButton aria-label="Search" onClick={refresh}><Search/></IconButton>
    </Tooltip>
  </div>
)

const ClearButton = ({ onClick }) => (
  <IconButton onClick={ onClick }>
    <Clear/>
  </IconButton>
)

class HeaderLine extends React.Component {
  constructor(props) {
    super(props)
    let searchQuery = window.location.pathname.split('/search/')[1] || ''
    searchQuery = unescape(searchQuery)
    this.state = { searchQuery }
  }

  get searchQuery() {
    return this.state.searchQuery
  }

  updateSearchQuery = (event) => {
    const searchQuery = event.target && event.target.value
    this.setState({ searchQuery })
  }

  submitSearch = (event) => {
    history.push(`/search/${this.searchQuery}`)
    event.preventDefault()
  }

  clearSearch = (event) => {
    this.setState({ searchQuery: '' })
    history.push(`/search`)
  }

	render() {
    const { classes, refresh, eventsContainer } = this.props

    return (
      <div>
        <AppBar position="static" color="default">
          <Toolbar className={classes.root}>
            <Typography
              variant="h6"
              color="inherit"
              className={classes.title}
              title={apiURLPrefix()}
            >Happening</Typography>
            <form className={classes.form} onSubmit={this.submitSearch}>
              <div style={{position: 'relative', display: 'inline-block', width: '100%' }}>
                <SearchPrefix refresh={this.submitSearch}/>
                <TextField
                  autoFocus={true}
                  id='searchbar'
                  label='Search'
                  value={this.searchQuery}
                  onChange={this.updateSearchQuery}
                  className={classes.formTextField}
                />
              </div>
            </form>
            <ClearButton onClick={this.clearSearch}/>
            <TimeMenu eventsContainer={eventsContainer} refresh={refresh}/>
            <RefreshButton refresh={refresh}/>
            <ChecksButton/>
            <LoginLogout/>
          </Toolbar>
        </AppBar>
      </div>
    )
  }
}

export default withStyles(styles)(HeaderLine)
