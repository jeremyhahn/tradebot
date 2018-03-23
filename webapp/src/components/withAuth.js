import React, { Component } from 'react';
import {BrowserRouter, browserHistory} from 'react-router-dom';
import { withRouter } from 'react-router'
import Login from 'app/components/Login';
import AuthService from 'app/components/AuthService';
import createHistory from 'history/createBrowserHistory'

export default function withAuth(AuthComponent) {

    var protocol = (window.location.protocol === "https:") ? "wss" : "ws";
    const Auth = new AuthService(window.location.protocol + '://localhost:8080/login');
    const history = createHistory();

    return withRouter(class AuthWrapped extends React.Component {

        constructor(props) {
            super(props);
            this.state = {
                user: null
            }
        }

        componentWillMount() {
            if (!Auth.loggedIn()) {
                history.replace('/login')
            }
            else {
                try {
                    const profile = Auth.getProfile()
                    this.setState({
                        user: Auth.getUser()
                    })
                }
                catch(err){
                  console.error('Unable to load JWT profile');
                    Auth.logout();
                    history.replace('/login')
                }
            }
        }

        render() {
          if (this.state.user) {
              return (
                  <AuthComponent history={this.history} user={this.state.user} />
              )
          }
          else {
              console.error('Authentication required!')
              return (
                  <Login history={this.history}/>
              )
          }
        }

    })
}
