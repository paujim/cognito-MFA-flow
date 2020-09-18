import React from 'react'

const authURL = process.env.REACT_APP_AUTH_URL
const localStorageKey = '__auth_provider_token__'

const getToken = async () => {
  // this is where we make a request to retrieve the user's token.
  return window.localStorage.getItem(localStorageKey)
}

const handleLoginResponse = (response) => {
  window.localStorage.setItem(localStorageKey, response.accessToken)
  return response.accessToken
}

const login = (username, password) => {
  return fetch_client('token/', "POST", { username, password })
    .then(handleLoginResponse)
    .catch(async (err) => {
      if (err.status) {
        let data = await err.json()
        throw data.message
      }
      console.log({ err })
      throw err
    })
}

const logout = async () => {
  window.localStorage.removeItem(localStorageKey)
}

const fetch_client = async (endpoint, method, data) => {
  const config = {
    method: method,
    body: data ? JSON.stringify(data) : undefined,
    headers: { 'Content-Type': data ? 'application/json' : undefined },
  }

  return fetch(`${authURL}/${endpoint}`, config)
    .then(response => {
      if (!response.ok) {
        throw response
      }
      return response.json()  //we only get here if there is no error
    })
}

const AuthContext = React.createContext()
AuthContext.displayName = 'AuthContext'

const useAuth = {
  getToken,
  login,
  logout,
}

export { useAuth, AuthContext }
