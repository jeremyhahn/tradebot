'use strict';

import React from 'react';

import AppBar from 'material-ui/AppBar';
import PropTypes from 'prop-types';
import Drawer from './Drawer';
import Toolbar from 'material-ui/Toolbar';
import { withRouter } from 'react-router-dom';
import { withStyles } from 'material-ui/styles';
import Typography from 'material-ui/Typography';
import IconButton from 'material-ui/IconButton';
import MenuIcon from 'material-ui-icons/Menu';

const styles = {
	root: {
		width: '100%',
	},
	flex: {
		flex: 1,
	},
	menuButton: {
		marginLeft: -12,
		marginRight: 20,
	},
};

class Header extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			drawer: false
		};
		this.handleDrawerToggle = this.handleDrawerToggle.bind(this);
		this.handleDrawerChange = this.handleDrawerChange.bind(this);
		this.handleTitleTap = this.handleTitleTap.bind(this);
		this.navigate = this.navigate.bind(this);
	}

	handleDrawerToggle() {
		this.setState({ drawer: !this.state.drawer });
	}

	handleDrawerChange(status) {
		this.setState({ drawer: status });
	}

	handleTitleTap() {
		this.props.history.push('/');
	}

	navigate(route) {
		return this.props.history.push(route)
	}

	render() {

		return (
			<div className="component--appbar">
				<AppBar>
					<Toolbar>
				    <IconButton color="inherit" aria-label="Menu" onClick={this.handleDrawerToggle}>
					    <MenuIcon  />
				    </IconButton>
				    <Typography type="title" color="inherit" >
					    { this.props.title }
				    </Typography>
				  </Toolbar>
				</AppBar>
				<Drawer open={this.state.drawer} change={this.handleDrawerToggle} navigate={this.navigate}/>
			</div>
		)
	}

}

export default withRouter(Header);
