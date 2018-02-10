import React from 'react';
import PropTypes from 'prop-types';
import LoginForm from 'app/components/LoginForm';

class Login extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      errors: {},
      user: {
        password: ''
      }
    };
    this.processForm = this.processForm.bind(this);
    this.changeUser = this.changeUser.bind(this);
  }

  processForm(event) {
    // prevent form submission event
    event.preventDefault();

    console.log('password:', this.state.user.password);
  }

  changeUser(event) {
    const field = event.target.name;
    const user = this.state.user;
    user[field] = event.target.value;
    this.setState({
      user
    });
  }

  render() {
    return (
      <LoginForm
        onSubmit={this.processForm}
        onChange={this.changeUser}
        errors={this.state.errors}
        user={this.state.user}
      />
    );
  }

}

export default Login;
