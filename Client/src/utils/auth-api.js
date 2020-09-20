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

const getUser = () => {
  let token = getToken()
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
  return client('token/', "POST", { username, password })
    .then(response => {
      return response.data
    })
    .catch(async (error) => {
      if (error.response) {
        // that falls out of the range of 2xx
        let data = await error.response.data
        throw data.message
      }
      throw "Unable fetch credentials"
    })
}

const changePassword = (username, password, session) => {
  return client('token/update', "POST", { username, password, session })
    .then(response => {
      return response.data
    })
    .catch(async (error) => {
      if (error.response) {
        // that falls out of the range of 2xx
        let data = await error.response.data
        throw data.message
      }
      throw "Unable update password"
    })
}

const logout = async () => {
  window.localStorage.removeItem(localStorageKey)
  window.location.reload();
}

const client = (endpoint, method, data) => {
  const config = {
    method: method,
    url: `${authURL}/${endpoint}`,
    data,
  }
  return axios(config)
}

const useAuthAPI = {
  setSession,
  getSession,
  getToken,
  setToken,
  getUser,
  login,
  changePassword,
  logout,
  isAuthenticated,
}

export { useAuthAPI }
