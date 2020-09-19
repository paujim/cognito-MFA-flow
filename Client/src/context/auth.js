import React from 'react'
import axios from 'axios'

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
  // this is where we make a request to retrieve the user's token.
  return window.localStorage.getItem(localStorageKey)
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
    .then( response => {
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
    .then( response => {
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
}

const client = (endpoint, method, data) => {
  const config = {
    method: method,
    url:`${authURL}/${endpoint}`,
    data,
  }
  return axios(config)
}

const AuthContext = React.createContext()
AuthContext.displayName = 'AuthContext'

const useAuth = {
  setSession,
  getSession,
  getToken,
  setToken,
  login,
  changePassword,
  logout,
  isAuthenticated,
}

export { useAuth, AuthContext }
