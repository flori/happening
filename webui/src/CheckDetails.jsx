import React from 'react'
import { Subscribe } from 'unstated'
import ChecksContainer from './ChecksContainer'
import Check from './CheckDetailed'

const checksContainer = new ChecksContainer()

export default class CheckDetails extends React.Component {
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
        <Subscribe to={[checksContainer]}>
          {cc => {
            const check = cc.checks.find(c => this.props.match.params.name === c.name)
            if (check) {
              return <Check refresh={refresh} update={update} key={check.id} {...check}/>
            } else {
              return null
            }
          }}
        </Subscribe>
      </>
    )
  }
}
