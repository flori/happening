import { Container } from 'unstated'
import { apiGetChecks } from './Api'

export default class ChecksContainer extends Container {
  state = {
    checks: [],
  }

  get checks() {
    return this.state.checks
  }

  load() {
    apiGetChecks(({ data: { data } }) => {
      this.setState({ checks: data || [] })
    })
  }
}

