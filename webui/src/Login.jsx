import React from 'react'
import PropTypes from 'prop-types'
import {
  Button,
  TextField,
  withStyles,
} from '@material-ui/core'
import { jwtLogin, getAuth, setAuth } from './Api'

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    flexDirection: 'column',
    alignItems: 'center',
  },
  textField: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    width: 200,
  },
  button: {
    marginTop: theme.spacing(2),
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    width: 200,
  }
})

class Login extends React.Component {
  state = {
    username: '',
    password: '',
  }

  handleChange = name => event => {
    this.setState({
      [name]: event.target.value,
    })
  }

  handleSubmit = (e) => {
    e.preventDefault()
    jwtLogin(this.state.username, this.state.password, this.props.to)
  }

  render() {
    if (getAuth()) {
      if (this.props.to) {
        window.location.href = '/?=' + this.props.to
      } else {
        window.location.href = '/'
      }
      return
    }
    const { classes } = this.props

    return <form className={classes.container} noValidate autoComplete="off">
      <TextField
        id="Username"
        label="Username"
        className={classes.textField}
        value={this.state.username}
        onChange={this.handleChange('username')}
        margin="normal"
        autoFocus={true}
      />
      <TextField
        id="Password"
        label="Password"
        className={classes.textField}
        value={this.state.password}
        onChange={this.handleChange('password')}
        margin="normal"
        type="password"
      />
      <Button size="medium" className={classes.button} variant="contained" color="primary" type="submit" onClick={this.handleSubmit}>
      Login
      </Button>
  </form>
  }
}

Login.propTypes = {
  classes: PropTypes.object.isRequired,
}

export default withStyles(styles)(Login)
