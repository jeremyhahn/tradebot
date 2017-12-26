'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import muiThemeable from 'material-ui/styles/muiThemeable';

import RaisedButton from 'material-ui/RaisedButton';

const container_style = {
	padding: 50,
	background: 'transparent',
	display: 'flex',
	flexDirection: 'column',
	justifyContent: 'center',
	marginTop: -100,
}

const EmptyExchangeList = (props) => {

	return (
		<div className="full-height-wrapper">
		<Paper style={ container_style } zDepth={1} rounded={false}>

			<h1 style={{ color: props.muiTheme.palette.textColor }}>Welcome!</h1>
			<h2 style={{ color: props.muiTheme.palette.textColor, marginTop: 20 }}>You have not added any exchanges yet. As soon as you do, your coins will appear here.</h2>

			<RaisedButton
				label="Add New Exchange"
				style={{ marginTop: 20 }}
				primary={true}
				onTouchTap={ props.openModal } />

		</Paper>
		</div>
	)
}

export default muiThemeable()(EmptyExchangeList);
