import axios from 'axios'
import Cookies from 'universal-cookie'

function handleUnauthorized(error) {
  if (!error.response) {
    return false
  }

  const { status } = error.response

  if (status === 401) {
    clearAuth()
    window.location.href = "/login"
    return true
  } else {
    return false
  }
}

export const apiInit = () => {
  axios.interceptors.response.use((response) => {
    return response
  }, (error) => {
    if (!handleUnauthorized(error)) {
      return Promise.reject(error)
    }
  })
}

export function getAuth() {
  const cookies = new Cookies()
  return cookies.get('auth')
}

export function clearAuth() {
  const cookies = new Cookies()
  cookies.remove('auth', { path: '/' })
}

export function setAuth(username, password) {
  const cookies = new Cookies()
  cookies.set('auth', `${username}:${password}`, {
    path: '/',
    expires:  new Date(new Date().getTime() + 14 * 86400 * 1000),
  })
}

export function getEnv(name, defaultName) {
  if (typeof(window.Env) !== 'undefined' && window.Env[name]) {
    return window.Env[name]
  } else {
    return defaultName
  }
}

export function apiURLPrefix() {
  return getEnv('HAPPENING_SERVER_URL', 'http://localhost:8080')
}

function buildApiURL(path) {
  let auth = getAuth()
  if (!auth) {
    window.location.href = "/login"
  }
  const [ username, password ] = auth.split(":", 2)
  return {
    username,
    password,
    url: new URL(apiURLPrefix() + path)
  }
}

function buildEventSearch(params) {
  let { query, seconds } = params

  if (!query) {
    query = ''
  } else {
    query = unescape(query)
  }

  const re    = /(id|name|output|hostname|command|success|ec):(\S+)/
  let m       = true
  let filters = {}
  while (m) {
    m = query.match(re)
    if (m) {
      filters[m[1]] = m[2]
      query = query.replace(new RegExp(`\\s*${m[1]}:\\S+\\s*`), ' ')
    }
  }

  query = query.replace(/(^\s+|\s+$)/g, '')

  let search
  if (query === '') {
    search = '?l=2000'
  } else {
    search = `?l=*&q=${query}`
  }

  if (seconds) {
    const s = Math.ceil(new Date().getTime() / 1000) - seconds
    search += `&s=${s}`
  }

  for (let f in filters) {
    search += `&f:${f}=${filters[f]}`
  }

  return search
}

export function apiGetEvents(params, block, handleError) {
  const path = '/api/v1/events'
  console.log(`Getting ${path} with ${JSON.stringify(params)}…`)
  let { url, username, password } = buildApiURL(path)
  url.search = buildEventSearch(params)
  axios.get(url, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiGetEvent(id, block, handleError) {
  const path = `/api/v1/event/${id}`
  console.log(`Getting ${path}…`)
  let { url, username, password } = buildApiURL(path)
  axios.get(url, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiGetChecks(block, handleError) {
  const path = '/api/v1/checks'
  console.log(`Getting ${path}…`)
  let { url, username, password } = buildApiURL(path)
  axios.get(url, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiGetCheckByName(name, block, handleError) {
  const path = `/api/v1/check/by-name/${name}`
  console.log(`Getting ${path}…`)
  let { url, username, password } = buildApiURL(path)
  axios.get(url, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiGetCheckById(id, block, handleError) {
  const path = `/api/v1/check/${id}`
  console.log(`Getting ${path}…`)
  let { url, username, password } = buildApiURL(path)
  axios.get(url, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiDeleteCheck({ id }, block, handleError) {
  const path = `/api/v1/check/${id}`
  console.log(`Deleting ${path}…`)
  let { url, username, password } = buildApiURL(path)
  axios.delete(url, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiPatchCheck(id, check, block, handleError) {
  const path = `/api/v1/check/${id}`
  console.log(`Patching ${path} with ${JSON.stringify(check)}…`)
  let { url, username, password } = buildApiURL(path)
  axios.patch(url, check, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiPutCheck(check, block, handleError) {
  const path = `/api/v1/check`
  console.log(`Putting ${path} with ${JSON.stringify(check)}…`)
  let { url, username, password } = buildApiURL(path)
  axios.put(url, check, { auth: { username, password } })
    .then(block)
    .catch(handleError)
}

export function apiStoreCheck(check, block, handleError) {
  apiGetCheckByName(check.name,
    (response) => { apiPatchCheck(response.data.data[0].id, check, block, handleError) },
    (error) => {
      const { response } = error
      if (response.status === 404) {
        apiPutCheck(check, block, handleError)
      } else {
        handleError(error)
      }
    }
  )
}
