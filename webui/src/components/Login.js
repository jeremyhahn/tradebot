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
    margin: '0 auto',
    width: '100px',
  },
  form: {
    marginLeft: '400',
    width:  '100px',
    height: '100px'
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    width: 200,
  },
  error: {
    color: 'red',
    padding: 25
  }
});

class Login extends Component {
    constructor() {
      super();
      this.state = {
        errors: null
      },
      this.handleChange = this.handleChange.bind(this);
      this.handleFormSubmit = this.handleFormSubmit.bind(this);
      this.Auth = new AuthService();
    }

    componentWillMount(){
      if(this.Auth.loggedIn())
          this.props.history.replace('/');
    }

    handleFormSubmit(e){
       e.preventDefault();
       this.Auth.login(this.state.username, this.state.password)
       .then(res =>{
          if(res.success) {
            this.props.history.replace('/portfolio');
          } else {
            this.setState({errors : res.error})
          }
       })
       .catch(err =>{
          console.error(err);
       })
    }

    handleChange(e){
        this.setState({
            [e.target.name]: e.target.value
        })
    }

    render() {

        const { classes } = this.props;

        return (
            <div className="center">
                <div className="card">
                    <h1>Welcome</h1>
                    {this.state.errors != null &&
                      <h3 className={classes.error}>{this.state.errors}</h3>
                    }
                    <form className="classes.form" noValidate autoComplete="off" onSubmit={this.handleFormSubmit}>

                        <TextField
                            id="username"
                            name="username"
                            label="Username"
                            type="username"
                            className={classes.textField}
                            onChange={this.handleChange}/>

                        <TextField
                            id="password"
                            name="password"
                            label="Password"
                            type="password"
                            className={classes.textField}
                            onChange={this.handleChange}/>


                        <br/>
                        <input className="form-submit" value="Login" type="submit"/>

                        <br/>
                        <Button className="form-submit" component={Link} to="/register">Register</Button>

                    </form>
                </div>
            </div>
        );
    }

}

export default withStyles(styles)(Login);
