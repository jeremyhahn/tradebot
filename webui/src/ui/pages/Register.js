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
    width: 200,
  }
});

class Register extends React.Component {

  constructor(props) {
		super(props);
    this.state = {
      open: true,
      password: "",
      confirm: "",
      passwordsMatch: false
    };
    this.handleRegister = this.handleRegister.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
    this.Auth = new AuthService();
  }

  handleRegister() {
    this.Auth.register(this.state.username, this.state.password)
    .then(res =>{
      this.props.history.replace('/');
    })
    .catch(err =>{
      console.error(err);
    })
  }

  handleCancel() {
    this.setState({open: false})
    this.props.history.push('/login');
  }

  handlePasswordChange(e) {
    const name = e.target.name
    const value = e.target.value
    this.setState({[name]: value},
        () => { this.validatePasswords(name, value)})
  }

  validatePasswords(name, value) {
    this.setState({passwordsMatch: this.state.password === this.state.confirm})
  }

  render() {

    const { classes } = this.props;

    return (
      <Dialog
        open={this.state.open}
        onClose={this.state.close}
        aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">Account Setup</DialogTitle>
        <DialogContent>
          <DialogContentText>
            <form className="classes.form" noValidate autoComplete="off">
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
                  error={!this.state.passwordsMatch}
                  onChange={(event) => this.handlePasswordChange(event)}/>
            </form>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleCancel} color="primary">Cancel</Button>
          <Button onClick={this.handleRegister} color="primary" disabled={!this.state.passwordsMatch}>Create</Button>
        </DialogActions>
      </Dialog>
    )
  }
};

Register.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Register);
