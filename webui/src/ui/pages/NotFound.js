'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';

const NotFound = (props) => {

	return (
		<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>
			<h1>Page Not Found.</h1>
			<p>You can try going back.</p>
		</Paper>
	)

}

export default NotFound;
