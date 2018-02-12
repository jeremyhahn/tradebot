import React, { Component } from 'react';
import AuthService from 'app/components/AuthService';

export default function withAuth(AuthComponent) {

var protocol = (loc.protocol === "https:") ? "wss" : "ws";

    const Auth = new AuthService(window.location.protocol + '://localhost:8080/login');

    return class AuthWrapped extends Component {
        constructor() {
            super();
            this.state = {
                user: null
            }
        }

        componentWillMount() {
            if (!Auth.loggedIn()) {
                this.props.history.replace('/login')
            }
            else {
                try {
                    const profile = Auth.getProfile()
                    this.setState({
                        user: profile
                    })
                }
                catch(err){
                    Auth.logout()
                    this.props.history.replace('/login')
                }
            }
        }

        render() {
          if (this.state.user) {
              return (
                  <AuthComponent history={this.props.history} user={this.state.user} />
              )
          }
          else {
              return null
          }
        }

    }
}
