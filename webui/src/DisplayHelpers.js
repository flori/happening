import { commandString } from './Event'
import moment from 'moment/moment'

const prettyBytes = require('pretty-bytes')

export const milli = 1000
export const micro = 1000000
export const nano = 1000000000

export function renderDate(date) {
  return moment(date).format()
}

export function renderDuration(durationNS) {
  if (isNaN(durationNS)) {
    return
  }
  const date = new Date(null)
  const s = durationNS / nano
  const d = Math.floor(s / 86400)
  date.setSeconds(s - d * 86400)
  const iso = date.toISOString().substr(11, 8)
  let rest = durationNS - Math.floor(s) * nano
  let result = ""
  if (d > 0) {
    result += d + "+" + iso
  } else {
    result += iso
  }
  rest = Math.round(rest / (nano / 1000))
  if (rest !== 0) {
    result += "." + rest
  }
  return result
}

export function parseDuration(durationString) {
  const re = /^(?:(\d+)?\+)?(\d{2}:\d{2}:\d{2})(\.\d+)?$/
  const m = durationString.match(re)
  if (durationString && m) {
    let seconds = 0
    if (m[1]) {
      seconds += parseInt(m[1], 10) * 86400
    }
    seconds += (new Date(`1970-01-01T${m[2]}Z`).getTime()) / 1000
    if (m[3]) {
      seconds += parseFloat(m[3])
    }
    return seconds * nano
  } else {
    return
  }
}

export function renderCommandResult({ success, exitCode, signal }) {
  let commandResult = ""
  const smiley = success ? "😀" : "😟"
  if (exitCode >= 0) {
    commandResult = `${smiley} ${exitCode}`
  } else {
    let signalName = signal || "n/a"
    commandResult = `${smiley} ${exitCode} (signal ${signalName})`
  }
  return commandResult
}

export function formatTooltip(event) {
    const { id, name, started, duration, load, hostname, user, cpuUsage, memoryUsage } = event
    let { command } = event
    if (command) {
      command = commandString(command)
      if (command.length > 40) {
        command = command.slice(0, 40) + '…'
      }
    } else {
      command = '–'
    }
    const startTime = Date.parse(started)
    const endTime = startTime + duration / micro
    return `
      <ul style="margin: 0.4em">
        <li>Name: <b>${name}</b></li>
        <li>Id: ${id}</li>
        <li>Started: ${renderDate(startTime)}</li>
        <li>Ended: ${renderDate(endTime)}</li>
        <li>Duration: ${renderDuration(duration)}</li>
        <li>Load: ${(100 * load).toFixed(2)}%</li>
        <li>CPU Usage: ${cpuUsage}</li>
        <li>Memory Usage: ${prettyBytes(memoryUsage)}</li>
        <li>Command: <tt>${command}</tt></li>
        <li>Success: ${renderCommandResult(event)}</li>
        <li>User: ${user}</li>
        <li>Hostname: ${hostname}</li>
      </ul>
    `
}
