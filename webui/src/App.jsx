import WebFont from 'webfontloader'
import React from 'react'
import qs from 'qs'

import {
  Router,
  Route,
  Redirect,
} from 'react-router-dom'

import HeaderLine from './HeaderLine'
import EventsList from './EventsList'
import ChecksList from './ChecksList'
import CheckDetails from './CheckDetails'
import Login from './Login'
import { Provider } from 'unstated'

import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'

import EventsContainer from './EventsContainer'

import { history } from './history'
import { apiInit, getAuth } from './Api'
apiInit()

WebFont.load({
  google: {
    families: ['Roboto:300,400,700', 'sans-serif', 'Material+Icons']
  }
})

const eventsContainer = new EventsContainer()

class Content extends React.Component {
  state = {
    update: 0
  }

  refresh = () => {
    this.setState((state) => ({ ...state, update: state.update + 1}))
  }

  render() {
    return <>
      <HeaderLine eventsContainer={eventsContainer} refresh={this.refresh}/>
      <Route path="/login" component={Login}/>
      <Route path="/search/:query"
        render={(props) => <EventsList
          update={this.state.update}
          eventsContainer={eventsContainer}
          refresh={this.refresh}
          {...props}
        />}
      />
      <Route exact path="/search/"
        render={(props) => <EventsList
          update={this.state.update}
          eventsContainer={eventsContainer}
          refresh={this.refresh}
          {...props}
        />}
      />
      <Route path="/check/:name"
        render={(props) => <CheckDetails refresh={this.refresh} update={this.state.update} {...props}/>}
      />
      <Route path="/checks"
        render={(props) => <ChecksList refresh={this.refresh} update={this.state.update} {...props}/>}
      />
      <Route exact path="/" render={( { location: { search } }) => {
        const params = qs.parse(search, { ignoreQueryPrefix: true })
        const p = params.p || '/search'
        return getAuth() ? <Redirect to={p}/> : <Login to={p}/>
      }}
      />
    </>
  }
}

const App = () => (
  <MuiThemeProvider theme={createMuiTheme({ typography: { useNextVariants: true }})}>
    <Router history={history}>
      <Provider>
        <Content/>
      </Provider>
    </Router>
  </MuiThemeProvider>
)

export default App
