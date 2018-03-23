import React, { Component } from 'react';
import { Link } from 'react-router-dom'
import AuthService from 'app/components/AuthService';
import Login from 'app/components/Login';

class Logout extends Component {

    constructor(props) {
      super(props);
      this.Auth = new AuthService();
    }

    componentWillMount() {
      this.Auth.logout();
    }

    render() {
        return (
          <Login history={this.props.history.replace('/login')} />
        );
    }

}

export default Logout;
