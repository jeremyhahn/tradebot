import decode from 'jwt-decode';
import axios from 'axios';

export default class AuthService {

    constructor(domain) {
        this.domain = domain || window.location.protocol + '//' +
            window.location.hostname + ':' + window.location.port + '/api/v1';
        this.fetch = this.fetch.bind(this);
    }

    login(username, password) {
        return this.fetch(`${this.domain}/login`, {
            method: 'POST',
            body: JSON.stringify({
              username: username,
              password: password
            })
        }).then(res => {
            this.setToken(res.token)
            return Promise.resolve(res);
        })
    }

    register(username, password) {
        return this.fetch(`${this.domain}/register`, {
            method: 'POST',
            body: JSON.stringify({
                username,
                password
            })
        }).then(res => {
            return Promise.resolve(res);
        })
    }

    syncTransactions() {
      return this.fetch(`${this.domain}/transactions/sync`, {
          method: 'GET'
      }).then(res => {
          return Promise.resolve(res)
      })
    }

    fetchTransactions(sync) {
      return this.fetch(`${this.domain}/transactions`, {
          method: 'GET'
      }).then(res => {
          return Promise.resolve(res)
      })
    }

    fetchOrderHistory() {
      return this.fetch(`${this.domain}/transactions/orderhistory`, {
          method: 'GET'
      }).then(res => {
          return Promise.resolve(res)
      })
    }

    exportTransactions() {
      return this.fetch(`${this.domain}/transactions/export`, {
          method: 'GET'
      }).then(res => {
          return Promise.resolve(res)
      })
    }

    importOrders(formData) {
      const config = {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'multipart/form-data'
          }
      }
      if(this.loggedIn()) {
          config.headers['Authorization'] = 'Bearer ' + this.getToken()
      }
      return axios.post(`${this.domain}/transactions/import`, formData, config)
    }

    updateCategory(id, formData) {
      const config = {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'multipart/form-data'
          }
      }
      if(this.loggedIn()) {
          config.headers['Authorization'] = 'Bearer ' + this.getToken()
      }
      return axios.put(`${this.domain}/transactions/${id}`, formData, config)
    }

    createExchange(formData) {
      const config = {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'multipart/form-data'
          }
      }
      if(this.loggedIn()) {
          config.headers['Authorization'] = 'Bearer ' + this.getToken()
      }
      return axios.post(`${this.domain}/user/exchange`, formData, config)
    }

    getExchangeNames() {
      return this.fetch(`${this.domain}/exchanges/names`, {
          method: 'GET'
      }).then(res => {
          return Promise.resolve(res)
      })
    }

    getUserExchanges() {
      return this.fetch(`${this.domain}/user/exchanges`, {
          method: 'GET'
      }).then(res => {
          return Promise.resolve(res)
      })
    }

    deleteUserExchange(exchangeName) {
      const config = {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'multipart/form-data'
          }
      }
      if(this.loggedIn()) {
          config.headers['Authorization'] = 'Bearer ' + this.getToken()
      }
      return axios.post(`${this.domain}/user/exchange/${exchangeName}`, null, config)
    }

    loggedIn() {
        const token = this.getToken()
        return !!token && !this.isTokenExpired(token)
    }

    isTokenExpired(token) {
        try {
            const decoded = decode(token);
            if (decoded.exp < Date.now() / 1000) {
              console.log('JWT expired');
              return true;
            }
            else
              return false;
        }
        catch (err) {
            return false;
        }
    }

    setToken(idToken) {
        console.log(idToken)
        localStorage.setItem('id_token', idToken)
        console.log(this.getProfile())
    }

    getToken() {
        return localStorage.getItem('id_token')
    }

    logout() {
        localStorage.removeItem('id_token');
    }

    getProfile() {
        var t = this.getToken()
        return t ? decode(t) : null
    }

    getUser() {
      var t = this.getProfile()
      return {
        id: t.user_id,
        username: t.username,
        local_currency: t.local_currency
      }
    }

    getExpiration() {
      var t = this.getProfile()
      return t["exp"]
    }

    fetch(url, options) {
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
        if (this.loggedIn()) {
            headers['Authorization'] = 'Bearer ' + this.getToken()
        }
        return fetch(url, {
            headers,
            ...options
        })
        .then(this._checkStatus)
        .then(response => response.json())
    }

    _checkStatus(response) {
        if (response.status == 200) {
            return response
        } else {
            var error = new Error(response.statusText)
            error.response = response
            throw error
        }
    }

}
