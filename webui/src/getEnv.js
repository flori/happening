export default function getEnv(name, defaultName) {
  if (name !== undefined) {
    const reactAppName = `REACT_APP_${name}`
    const reactAppValue = process.env[reactAppName]
    if (reactAppValue !== undefined) {
      return reactAppValue
    }
  }
  if (typeof(window.Env) !== 'undefined' && window.Env[name]) {
    return window.Env[name]
  }
  return defaultName
}
