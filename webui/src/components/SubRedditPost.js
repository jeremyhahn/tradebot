'use strict';

import React from 'react';

import { List, ListItem } from 'material-ui/List';
import { withRouter } from 'react-router-dom';
import Divider from 'material-ui/Divider';


const SubRedditPost = (props) => {

	const { data, history } = props;

	// when clicked run this function
	const handleTap = () => {
		return history.push(data.url);
	}

	return (
		<div>
		<ListItem
			primaryText={ data.title }
			secondaryText={ data.description }
			// onTouchTap={ handleTap }
		/>
		<Divider />
		</div>
	)
}

export default withRouter( SubRedditPost );

