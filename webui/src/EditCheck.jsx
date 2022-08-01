import React from 'react'
import {
  Icon,
  IconButton,
  ListItemIcon,
} from '@material-ui/core'
import { apiGetCheckByNameInContext, apiPatchCheck, apiPutCheck } from './Api'
import Alert from './Alert'
import EditCheckDialog from './EditCheckDialog'

export default class EditCheck extends React.Component {
  state = {
    check: null,
  }

  loadCheck() {
    if (this.props.name && this.props.context) {
      apiGetCheckByNameInContext(
        this.props.name,
        this.props.context,
        ({ data: { data } }) => {
          this.setState({ check: data[0] })
        },
        (error) => {
          if (error.status === 404) {
            this.setState({
              check: {
                name: this.props.name,
              }
            })
          }
        }
      )
    } else {
      this.setState({
        check: {
          name: "",
          context: "default"
        }
      })
    }
  }

  setMessage = (error) => { this.setState({ message: error.message }) }

  handleClickOpen = () => {
    this.loadCheck()
  }

  handleClose = () => {
    this.setState({ check: null });
  }

  handleCloseSave = check => {
    switch (this.props.action) {
      case "edit":
        apiPatchCheck(
          check.id,
          check,
          this.props.refresh,
          this.setMessage
        )
        break
      case "add":
        apiPutCheck(
          check,
          this.props.refresh,
          this.setMessage
        )
        break
      default:
    }
    this.setState({ check: null });
  }

  render() {
    const {
      action,
      name,
      context,
    } = this.props
    const check = this.state.check
    return (
      <>
        <Alert variant="error" message={this.state.message} duration={6000} reset={() => { this.setState({ message: null }) } }/>
        <ListItemIcon>
          <IconButton aria-label={action} onClick={this.handleClickOpen}>
            <Icon>{action}_box</Icon>
          </IconButton>
        </ListItemIcon>
        {this.state.check && <EditCheckDialog
          action={action}
          name={name}
          context={context}
          disabled={check.disabled}
          failures={check.failures}
          allowed_failures={check.allowed_failures}
          success={check.success}
          healthy={check.healthy}
          last_ping_at={check.last_ping_at}
          onOpen={this.handleClickOpen}
          onClose={this.handleClose}
          onCloseSave={this.handleCloseSave}
          open={!!this.state.check}
          id={check.id}
          period={check.period}
          refresh={check.refresh}
        />}
      </>
    )
  }
}
