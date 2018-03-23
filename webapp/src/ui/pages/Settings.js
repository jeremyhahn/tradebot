'use strict';

import React from 'react';
import PropTypes from 'prop-types';
import Paper from 'material-ui/Paper';
import Loading from 'app/components/Loading';
import { withStyles } from 'material-ui/styles';
import AuthService from 'app/components/AuthService';
import withAuth from 'app/components/withAuth';

const styles = theme => ({
  root: {
    flex: 1,
    paddingLeft: '1%',
    width: '99%',
    marginTop: '68px'
  },
  table: {
    width: '100%'
  },
  tableWrapper: {
    overflowX: 'auto',
  }
});

class Settings extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			loading: false
		}
		this.Auth = new AuthService();
	}

	componentDidMount() {

	}

	render() {

		return (
			<Paper style={{ padding: '20px', marginTop: '60px' }}>
				<h2>Settings</h2>

				{this.state.loading &&
					<Loading text="Loading settings..." />
				}
			</Paper>
		)

	}
}

export default withAuth(withStyles(styles)(Settings));
