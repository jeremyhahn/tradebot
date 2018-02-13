import React, { Component } from 'react';
import {BrowserRouter} from 'react-router';
import AuthService from 'app/components/AuthService';
import { createHashHistory } from 'history'

export default function withAuth(AuthComponent) {

    var protocol = (window.location.protocol === "https:") ? "wss" : "ws";
    const Auth = new AuthService(window.location.protocol + '://localhost:8080/login');
    const history = createHashHistory();

    return class AuthWrapped extends Component {

        constructor(p) {
            super();
            this.state = {
                user: null,
                tempUser: {
                  id: 1,
                  username: "test",
                  local_currency: "USD"
                }
            }
        }

        componentWillMount() {
            if (!Auth.loggedIn()) {
                //history.replace('/login')
            }
            else {
                try {
                    const profile = Auth.getProfile()
                    this.setState({
                        user: profile
                    })
                }
                catch(err){
                    //Auth.logout();
                    console.log('Unable to load JWT profile');
                    //history.replace('/login')
                }
            }
        }

        render() {
          if (this.state.user) {
              return (
                  <AuthComponent history={history} user={this.state.user} />
              )
          }
          else {
              console.error('withAuth Authentication failed!')
              return (
                  <AuthComponent user={this.state.tempUser} />
              )
          }
        }

    }
}
