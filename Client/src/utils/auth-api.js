import axios from 'axios'
import jwt from 'jwt-decode'

const authURL = process.env.REACT_APP_AUTH_URL
const localStorageKey = '__auth_provider_token__'
const localStorageSession = '__auth_provider_session__'

const getSession = () => {
  return window.localStorage.getItem(localStorageSession)
}
const setSession = (session) => {
  return window.localStorage.setItem(localStorageSession, session)
}

const getToken = () => {
  return window.localStorage.getItem(localStorageKey)
}

const decodeUser = (token) => {
  return token ? jwt(token) : { username: "Guest" }
}

const setToken = (token) => {
  return window.localStorage.setItem(localStorageKey, token)
}

const isAuthenticated = () => {
  let token = getToken()
  return token !== "" && token !== undefined && token !== null
}

const login = (username, password) => {
  return client('token/', "POST", { username, password }, "Unable fetch credentials")
}

const changePassword = (username, password, session) => {
  return client('token/update', "POST", { username, password, session }, "Unable update password")
}

const verifyCode = (username, code, session) => {
  return client('token/code', "POST", { username, code, session }, "Unable verify code")
}

const registerMFA = (accessToken) => {
  return client('mfa/register', "POST", { accessToken }, "Unable register MFA")
}

const verifyMFA = (accessToken, code) => {
  return client('mfa/verify', "POST", { accessToken, code, deviceName:"Galaxy" }, "Unable verify code")
}

const logout = () => {
  window.localStorage.removeItem(localStorageKey)
  window.location.reload();
}

const client = (endpoint, method, data, defaultErrorMessage) => {
  const config = {
    method: method,
    url: `${authURL}/${endpoint}`,
    data,
  }
  return axios(config)
    .then(response => response.data)
    .catch(async (error) => {
      if (error.response) {
        // that falls out of the range of 2xx
        let data = await error.response.data
        throw new Error(data.message)
      }
      throw new Error(defaultErrorMessage)
    })
}

const useAuthAPI = {
  setSession,
  getSession,
  getToken,
  setToken,
  decodeUser,
  login,
  changePassword,
  logout,
  isAuthenticated,
  registerMFA,
  verifyMFA,
  verifyCode,
}

export { useAuthAPI }
