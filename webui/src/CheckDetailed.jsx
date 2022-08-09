import React from 'react'
import {
  Avatar,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  ListItemIcon,
  Typography,
} from '@material-ui/core'
import SearchButton from './SearchButton'
import ConfirmDeleteCheck from './ConfirmDeleteCheck'
import EditCheck from './EditCheck'
import { checkState, CheckStateIcon } from './checkState'
import { renderDuration, renderDate, micro } from './DisplayHelpers'

export default class CheckDetailed extends React.Component {
  render() {
    const {
      id,
      name,
      context,
      healthy,
      success,
      failures,
      allowed_failures,
      period,
      last_ping_at,
      disabled,
      refresh,
    } = this.props

    const { stateDescription } = checkState({ healthy, success, disabled })

    const shouldRepeat = "every " + renderDuration(period)

    let noRepeat
    if (!disabled && !healthy && success) {
      const lastChance = new Date(new Date(last_ping_at).getTime() + period / micro)
      noRepeat = ", but didn't, should have repeated before " + renderDate(lastChance)
    }

    let title = success ? 'healthy' : 'unhealthy'

    const failureText = `${failures}/${allowed_failures} failed`
    if (allowed_failures > 0) {
      title += ` ${failureText}`
    }

    return (
      <>
        <Card>
          <CardHeader avatar={
            <Avatar title={title} aria-label={title}>
              <CheckStateIcon {...this.props}/>
            </Avatar>}
            title={`Name: ${name}`}
            subheader={`State: ${stateDescription}`}
          />
          <CardContent>
            <Typography align="left" variant="body2">Check is in state <strong>{stateDescription}</strong> now</Typography>
            <Typography align="left" variant="body2">Last run at {renderDate(last_ping_at)} was <strong>{success ? 'success' : 'failure'}</strong></Typography>
            <Typography align="left" variant="body2">Allowed to fail? <strong>{failureText}</strong></Typography>
            <Typography align="left" variant="body2">Should repeat {shouldRepeat}{noRepeat}</Typography>
          </CardContent>
          <CardActions>
            <ListItemIcon>
              <SearchButton name={name} context={context}/>
            </ListItemIcon>
            <EditCheck action="edit" name={name} context={context} refresh={refresh}/>
            <ConfirmDeleteCheck name={name} id={id} refresh={refresh}/>
          </CardActions>
        </Card>
    </>
    )
  }
}
