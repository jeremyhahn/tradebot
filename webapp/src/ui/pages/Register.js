import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import { withRouter } from 'react-router-dom';
import { Link } from 'react-router';
import Card, { CardActions, CardContent } from 'material-ui/Card';
import Button from 'material-ui/Button';
import { FormControl, FormHelperText } from 'material-ui/Form';
import TextField from 'material-ui/TextField';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';
import AuthService from 'app/components/AuthService';

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
  },
  error: {
    color: 'red',
    padding: 25
  },
  createButton: {
    opacity: 0.6
  }
});

class Register extends React.Component {

  constructor(props) {
		super(props);
    this.state = {
      password: "",
      confirm: "",
      passwordsMatch: false,
      errors: null,
      opacity: 0.6,
      disabled: true
    };
    this.handleRegister = this.handleRegister.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
    this.Auth = new AuthService();
  }

  handleRegister(e) {
    e.preventDefault();
    this.Auth.register(this.state.username, this.state.password)
    .then(res => {
      if(res.error) {
        this.setState({errors: res.error})
        return false;
      }
      this.props.history.replace('/');
    })
    .catch(err =>{
      console.error(err);
    })
  }

  handleCancel() {
    this.props.history.push('/login');
  }

  handlePasswordChange(e) {
    const name = e.target.name
    const value = e.target.value
    this.setState({[name]: value},
        () => { this.validatePasswords(name, value)})
  }

  validatePasswords(name, value) {
    if(this.state.password === this.state.confirm) {
      this.setState({
        passwordsMatch: true,
        opacity: 1,
        disabled: false})
    } else {
      this.setState({
        passwordsMatch: false,
        opacity: 0.6,
        disabled: true})
    }
    /*
    var strongRegex = new RegExp("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])(?=.{8,})");
    //var mediumRegex = new RegExp("^(((?=.*[a-z])(?=.*[A-Z]))|((?=.*[a-z])(?=.*[0-9]))|((?=.*[A-Z])(?=.*[0-9])))(?=.{6,})");
    if(!strongRegex.test(this.state.password)) {
      this.setState({
        errors: "Password not strong enough"
    }
    */
  }

  render() {

    const { classes } = this.props;

    return (
      <div className="center">
          <div className="card">
              <h1>Account Setup</h1>
              {this.state.errors != null &&
                <h3 className={classes.error}>{this.state.errors}</h3>
              }
              <form className="classes.form" noValidate autoComplete="off" onSubmit={((event) => this.handleRegister(event))}>
                <FormControl fullWidth className={classes.formControl}>
                  <TextField
                      id="username"
                      name="username"
                      label="Username"
                      type="username"
                      placeholder="(Optional)"
                      className={classes.textField}
                      value={this.state.username}
                      onChange={(event) => this.setState({[event.target.name]: event.target.value})}/>
                </FormControl>
                <TextField
                    id="password"
                    name="password"
                    label="Password"
                    type="password"
                    className={classes.textField}
                    value={this.state.password}
                    onChange={(event) => this.handlePasswordChange(event)}/>
                <TextField
                    id="confirm"
                    name="confirm"
                    label="Confirm"
                    type="password"
                    className={classes.textField}
                    value={this.state.confirm}
                    error={!this.state.passwordsMatch && this.state.password != ""}
                    onChange={(event) => this.handlePasswordChange(event)}/>

                <br/>
                <input className="form-submit" value="Create Account" type="submit"
                   style={{opacity: this.state.opacity}} disabled={this.state.disabled}/>

                <br/>
                <Button className="form-submit" onClick={this.handleCancel}>Cancel</Button>
              </form>

          </div>
      </div>
    )
  }
};

Register.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Register);
