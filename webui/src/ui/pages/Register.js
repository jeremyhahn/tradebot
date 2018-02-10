import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import { withRouter } from 'react-router-dom';
import { Link } from 'react-router';
import Card, { CardActions, CardContent } from 'material-ui/Card';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';

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
      open: true
    };
    this.handleClose = this.handleClose.bind(this);
  }

  handleChange(name) {
    console.log(name)
  }

  handleClose() {
    this.setState({open: false})
  }

  handleCreate() {
    console.log('create!')
  }

  render() {

    const { classes } = this.props;

    return (
      <Dialog
        open={this.state.open}
        onClose={this.state.close}
        aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">New Account</DialogTitle>
        <DialogContent>
          <DialogContentText>
            <form className="classes.form" noValidate autoComplete="off">
              <TextField
                  id="password"
                  label="Password"
                  type="password"
                  className={classes.textField}
                  value={this.state.name}
                  onChange={this.handleChange('password')}/>
              <TextField
                  id="confirm"
                  label="Confirm"
                  type="password"
                  className={classes.textField}
                  value={this.state.name}
                  onChange={this.handleChange('confirm')}/>
            </form>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleCreate} color="primary">Create</Button>
          <Button onClick={this.handleClose} color="primary">Cancel</Button>
        </DialogActions>
      </Dialog>
    )
  }
};

Register.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Register);
