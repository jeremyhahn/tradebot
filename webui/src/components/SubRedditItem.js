'use strict';

import React from 'react';

import { List, ListItem } from 'material-ui/List';
import { withRouter } from 'react-router-dom';


const SubRedditItem = (props) => {

	const { data, history } = props;

	// when clicked run this function
	const handleTap = () => {
		return history.push(data.url);
	}

	return (
		<ListItem
			primaryText={ data.title }
			secondaryText={ data.description }
			onTouchTap={ handleTap }
		/>
	)
}

export default withRouter( SubRedditItem );

