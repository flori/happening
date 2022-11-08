import React from 'react'
import { Link } from 'react-router-dom'
import LinkIcon from '@material-ui/icons/Link'
import {
  Chip,
  Divider,
  IconButton,
  ListItem,
  ListItemIcon,
  ListItemSecondaryAction,
  ListItemText,
} from '@material-ui/core'
import SearchButton from './SearchButton'
import ConfirmDeleteCheck from './ConfirmDeleteCheck'
import EditCheck from './EditCheck'
import { CheckStateAvatar } from './checkState'
import { renderDuration, renderDate, micro } from './DisplayHelpers'

export default class Check extends React.Component {
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

    const primaryText = name + " (every " + renderDuration(period) + ")"
    let secondaryText = "last pinged at " + renderDate(last_ping_at)

    if (!disabled && !healthy && success) {
      const lastChance = new Date(new Date(last_ping_at).getTime() + period / micro)
      secondaryText += ", should have repeated before " + renderDate(lastChance)
    }

    let title = success ? 'healthy' : 'unhealthy'

    if (allowed_failures > 0) {
      title += ` ${failures}/${allowed_failures} failed`
    }
    return (
      <>
        <ListItem>
          <CheckStateAvatar {...this.props} refresh={refresh}/>
          <ListItemText primary={primaryText} secondary={secondaryText}/>
          <ListItemSecondaryAction style={{ display: 'flex', flexFlow: 'row wrap', flexDirection: 'row' }}>
            <ListItemIcon>
              <Chip label={context} color="secondary"/>
            </ListItemIcon>
            <ListItemIcon>
              <IconButton title={title} aria-label={title} component={Link} to={`/check/${name}`}>
                <LinkIcon/>
              </IconButton>
            </ListItemIcon>
            <ListItemIcon>
              <SearchButton name={name} context={context}/>
            </ListItemIcon>
            <EditCheck action="edit" name={name} context={context} refresh={refresh}/>
            <ConfirmDeleteCheck name={name} id={id} refresh={refresh}/>
          </ListItemSecondaryAction>
        </ListItem>
      <Divider/>
    </>
    )
  }
}
