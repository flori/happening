import axios from 'axios'
import Cookies from 'universal-cookie'
import getEnv from './getEnv'

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

export function jwtLogin(username, password, to) {
  const authPair = { username, password }
  axios.post(
    apiURLPrefix() + "/jwt", authPair,
    { headers: { "Content-Type": "application/json" } })
    .then((data) => {
      const token = data.data.token
      setAuth(token)
      if (to) {
        window.location.href = '/?=' + to
      } else {
        window.location.href = '/'
      }
    })
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

export function setAuth(token) {
  const cookies = new Cookies()
  cookies.set('auth', token, {
    path: '/',
    expires:  new Date(new Date().getTime() + 7 * 86400 * 1000), // 1 week
  })
}

export function apiURLPrefix() {
  return getEnv('HAPPENING_SERVER_URL', 'http://localhost:8080')
}

function buildApiURL(path) {
  const token = getAuth()
  if (!token) {
    window.location.href = "/login"
  }
  return {
    token,
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

  const re    = /(id|name|context|output|hostname|command|success|ec):(\S+)/
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
  const { url, token } = buildApiURL(path)
  url.search = buildEventSearch(params)
  axios.get(url, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiGetEvent(id, block, handleError) {
  const path = `/api/v1/event/${id}`
  console.log(`Getting ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.get(url, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiPatchMailEvent(id, block, handleError) {
  const path = `/api/v1/event/${id}/mail`
  console.log(`Mailing ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.patch(url, {}, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiGetChecks(block, handleError) {
  const path = '/api/v1/checks'
  console.log(`Getting ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.get(url, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiGetCheckByNameInContext(name, context, block, handleError) {
  const path = `/api/v1/check/by-name/${name}/in-context/${context}`
  console.log(`Getting ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.get(url, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiGetCheckById(id, block, handleError) {
  const path = `/api/v1/check/${id}`
  console.log(`Getting ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.get(url, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiDeleteCheck({ id }, block, handleError) {
  const path = `/api/v1/check/${id}`
  console.log(`Deleting ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.delete(url, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiPatchCheck(id, check, block, handleError) {
  const path = `/api/v1/check/${id}`
  console.log(`Patching ${path} with ${JSON.stringify(check)}…`)
  const { url, token } = buildApiURL(path)
  axios.patch(url, check, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiResetCheck(id, block, handleError) {
  const path = `/api/v1/check/${id}/reset`
  console.log(`Resetting ${path}…`)
  const { url, token } = buildApiURL(path)
  axios.patch(url, {}, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}


export function apiPutCheck(check, block, handleError) {
  const path = `/api/v1/check`
  console.log(`Putting ${path} with ${JSON.stringify(check)}…`)
  const { url, token } = buildApiURL(path)
  axios.put(url, check, { headers: { "Authorization": `Bearer ${token}` } })
    .then(block)
    .catch(handleError)
}

export function apiStoreCheck(check, block, handleError) {
  apiGetCheckByNameInContext(check.name, check.context,
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
