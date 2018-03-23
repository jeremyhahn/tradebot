'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import Button from 'material-ui/Button';

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

			<h1 >Welcome!</h1>
			<h2 style={{ marginTop: 20 }}>You have not added any exchanges yet. As soon as you do, your coins will appear here.</h2>

			<Button
				label="Add New Exchange"
				style={{ marginTop: 20 }}
				primary={true}
				onTouchTap={ props.openModal } />

		</Paper>
		</div>
	)
}

export default EmptyExchangeList;
