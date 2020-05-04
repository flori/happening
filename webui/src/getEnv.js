export default function getEnv(name, defaultName) {
  if (name !== null && process.env[name] !== null) {
    return process.env[`REACT_APP_${name}`]
  }
  if (typeof(window.Env) !== 'undefined' && window.Env[name]) {
    return window.Env[name]
  }
  return defaultName
}
