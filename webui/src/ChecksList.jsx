import React from 'react'
import {
  List,
  ListItem,
  ListItemText,
} from '@material-ui/core'
import EditCheck from './EditCheck'
import { Subscribe } from 'unstated'
import ChecksContainer from './ChecksContainer'
import Check from './Check'

const checksContainer = new ChecksContainer()

export default class ChecksList extends React.Component {
  componentDidMount() {
    checksContainer.load()
  }

  componentDidUpdate() {
    checksContainer.load()
  }

  render() {
    const { refresh, update } = this.props
    return (
      <>
        <List>
          <ListItem>
            <ListItemText primary="Health" secondary="Checks"/>
            <EditCheck action="add" refresh={refresh} update={update}/>
          </ListItem>
          <Subscribe to={[checksContainer]}>
          {cc => cc.checks.map( (c) => <Check refresh={refresh} update={update} key={c.id} {...c}/> )}
          </Subscribe>
        </List>
      </>
    )
  }
}
