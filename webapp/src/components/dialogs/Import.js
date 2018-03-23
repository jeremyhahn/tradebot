import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import Button from 'material-ui/Button';
import { FormControl } from 'material-ui/Form';
import TextField from 'material-ui/TextField';
import Dialog, {
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from 'material-ui/Dialog';
import { InputLabel } from 'material-ui/Input';
import Select from 'material-ui/Select';
import { MenuItem } from 'material-ui/Menu';
import AuthService from 'app/components/AuthService';

const styles = theme => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  formControl: {
    margin: theme.spacing.unit,
    minWidth: 120,
  },
  selectEmpty: {
    marginTop: theme.spacing.unit * 2,
  },
});

class ImportDialog extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      file: null,
      exchange: "",
      exchanges: []
    }
    this.onFormSubmit = this.onFormSubmit.bind(this)
    this.onChange = this.onChange.bind(this)
    this.Auth = new AuthService();
  }

  componentDidMount() {
    this.Auth.getExchangeNames()
    .then((response) => {
      console.log(response)
      if(response.success) {
        this.setState({exchanges: response.payload})
      }
    })
  }

  handleChange = (event, value) => {
    this.setState({[event.target.name]: event.target.value});
  };

  onFormSubmit(e) {
    e.preventDefault()
    var _this = this
    const formData = new FormData();
    formData.append('csv', this.state.file)
    formData.append('exchange', this.state.exchange)
    this.Auth.importOrders(formData)
      .then((response) => {
        if(response.data.success) {
          _this.props.onClose()
          _this.props.addData(response.data.payload)
        }
    })
  }

  onChange(e) {
    this.setState({file: e.target.files[0] })
  }

  render() {

    const { classes } = this.props;

    return (
        <Dialog
          open={this.props.open}
          onClose={this.props.onClose}
          aria-labelledby="form-dialog-title">
          <DialogTitle id="form-dialog-title">Import CSV</DialogTitle>
          <form onSubmit={this.onFormSubmit}>
            <DialogContent>
                <FormControl className={classes.formControl} fullWidth={true}>
                <InputLabel htmlFor="exchange">Exchange</InputLabel>
                  <Select
                    value={this.state.exchange}
                    onChange={this.handleChange}
                    inputProps={{
                      name: 'exchange',
                      id: 'exchange',
                    }}>
                  <MenuItem value=""><em>None</em></MenuItem>
                  { this.state.exchanges.map( exchange =>
                    <MenuItem key={exchange} value={exchange}>{exchange}</MenuItem>
                  )}
                  </Select>
                </FormControl>
                <FormControl className={classes.formControl} fullWidth={true}>
                  <input type="file" name="csv" onChange={this.onChange} />
                </FormControl>
            </DialogContent>
            <DialogActions>
              <Button onClick={this.props.onClose} color="primary">Cancel</Button>
              <Button type="submit" label="submit" color="primary">Upload</Button>
            </DialogActions>
          </form>
        </Dialog>
    );
  }
}

export default withStyles(styles)(ImportDialog);
