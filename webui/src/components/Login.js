import React, { Component } from 'react';
import { Link } from 'react-router-dom'
import { withStyles } from 'material-ui/styles';
import Button from 'material-ui/Button';
import { FormControl, FormHelperText } from 'material-ui/Form';
import TextField from 'material-ui/TextField';
import AuthService from 'app/components/AuthService';
import 'app/css/login.css';

const styles = theme => ({
  container: {
    margin: '0 auto'
  },
  form: {
    marginLeft: '400',
    width:  '100px',
    height: '100px'
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  },
  error: {
    color: 'red',
    textAlign: 'center',
    marginTop: '25'
  }
});

class Login extends Component {

    constructor(props) {
      super(props);
      this.state = {
        password: "",
        opacity: 0.6,
        disabled: true,
        errors: ""
      },
      this.handleChange = this.handleChange.bind(this);
      this.handleFormSubmit = this.handleFormSubmit.bind(this);
      this.Auth = new AuthService();
    }

    componentWillMount(){
      if(this.Auth.loggedIn())
          this.props.history.replace('/orders');
    }

    handleFormSubmit(e){
       e.preventDefault();
       this.Auth.login(this.state.username, this.state.password)
       .then(res => {
         console.log(res.token.length)
          if(res.token.length) {
console.log('navigating to /orders')
            //this.props.history.replace('/orders');
            location.href = '/orders'

          } else {
            this.setState({errors : res.error})
          }
       })
       .catch(err => {
          if(err.response) {
            console.log(err.response)
            console.log(err.statusText)
            console.log(err)
          }
       })
    }

    handleChange(e){
        this.setState({
            [e.target.name]: e.target.value
        })
        if(this.state.password.length > 0) {
          this.setState({
            opacity: 1,
            disabled: false})
        } else {
          this.setState({
            opacity: 0.6,
            disabled: true
          })
        }
    }

    render() {

        const { classes } = this.props;

        return (
          <div className={classes.container}>
            <div className="center">
                <div className="card">
                    <h1>Welcome</h1>
                    {this.state.errors != "" &&
                      <h3 className={classes.error}>{this.state.errors}</h3>
                    }
                    <form className="classes.form" noValidate autoComplete="off" onSubmit={this.handleFormSubmit}>

                        <TextField
                            id="username"
                            name="username"
                            label="Username"
                            type="username"
                            value={this.state.username}
                            className={classes.textField}
                            onChange={this.handleChange}/>

                        <TextField
                            id="password"
                            name="password"
                            label="Password"
                            type="password"
                            value={this.state.password}
                            className={classes.textField}
                            onChange={this.handleChange}/>

                        <br/>
                        <input className="form-submit" value="Login" type="submit"
                           style={{opacity: this.state.opacity}} disabled={this.state.disabled}/>

                        <br/>
                        <Button className="form-submit" component={Link} to="/register">Register</Button>

                    </form>
                </div>
            </div>
        </div>
        );
    }

}

export default withStyles(styles)(Login);
