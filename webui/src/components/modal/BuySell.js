import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import Tabs, { Tab } from 'material-ui/Tabs';
import Button from 'material-ui/Button';
import Loading from 'app/components/Loading';
import Typography from 'material-ui/Typography';
import Dialog, { DialogTitle } from 'material-ui/Dialog';

function TabContainer(props) {
  return (
    <Typography component="div" style={{ padding: 8 * 3 }}>
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

const styles = theme => ({
  root: {
    flexGrow: 1,
    marginTop: theme.spacing.unit * 3,
    backgroundColor: theme.palette.background.paper,
  },
});

class BuySellModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      value: 0
    };
  }

  handleChange = (event, value) => {
    this.setState({ value });
  };

  handleClose = () => {
    this.props.onClose(this.props.selectedValue);
  };

  render() {

    const { classes } = this.props;
    const { value } = this.state;

		const actions = [
			<Button
				label="Cancel"
				primary={ true }
				onTouchTap={ this.props.close }
			/>,
			<Button
				label="Submit"
				primary={ true }
				disabled={ ! this.state.title || !this.state.url }
				onTouchTap={ this.submit }
			/>,
		];

    return (

        <Dialog
  				title="Sell"
  				actions={ actions }
  				open={ this.props.open }
          onClose={this.handleClose} >

          <DialogTitle id="buy-sell-title">Buy & Sell</DialogTitle>

  				{ this.state.processing &&
  					<div>
  						<Loading />
  					</div>
  				}

  				{ ! this.state.processing &&
            <Tabs value={value} onChange={this.handleChange}>
              <Tab label="Item One" />
              <Tab label="Item Two" />
              <Tab label="Item Three" href="#basic-tabs" />
            </Tabs>
  			  }

          {value === 0 && <TabContainer>Item One</TabContainer>}
          {value === 1 && <TabContainer>Item Two</TabContainer>}
          {value === 2 && <TabContainer>Item Three</TabContainer>}

  			</Dialog>
    );
  }
}

BuySellModal.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(BuySellModal);
