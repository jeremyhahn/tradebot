import decode from 'jwt-decode';

export default class AuthService {

    constructor(domain) {
        this.domain = domain || window.location.protocol + '//localhost:8080/api/v1';
        this.fetch = this.fetch.bind(this);
        /*
        this.register = this.register.bind(this);
        this.login = this.login.bind(this);
        this.logout = this.login.bind(this);
        this.setToken = this.setToken.bind(this);
        this.getToken = this.getToken.bind(this);
        this.getProfile = this.getProfile.bind(this);
        this.getUser = this.getUser.bind(this);
        this.getExpiration = this.getExpiration.bind(this);
        this.loggedIn = this.loggedIn.bind(this);
        this.isTokenExpired = this.isTokenExpired.bind(this);
        */
    }

    login(username, password) {
        return this.fetch(`${this.domain}/login`, {
            method: 'POST',
            body: JSON.stringify({
              username: username,
              password: password
            })
        }).then(res => {
            if(res.token.length) {
              this.setToken(res.token)
            }
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

    loggedIn() {
        const token = this.getToken()
        return !!token && !this.isTokenExpired(token)
    }

    isTokenExpired(token) {
        try {
            const decoded = decode(token);
            if (decoded.exp < Date.now() / 1000) {
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
        localStorage.setItem('id_token', idToken)
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
