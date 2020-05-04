export default function getEnv(name, defaultName) {
  if (name !== null) {
    const reactAppName = `REACT_APP_${name}`
    if (process.env[reactAppName] !== null) {
      return process.env[reactAppName]
    }
  }
  if (typeof(window.Env) !== 'undefined' && window.Env[name]) {
    return window.Env[name]
  }
  return defaultName
}
