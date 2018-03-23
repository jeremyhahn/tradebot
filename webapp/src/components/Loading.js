'use strict';

import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import { CircularProgress } from 'material-ui/Progress';

const styles = theme => ({
	container: {
		textAlign: 'center'
	},
	refresh: {
		display: 'inline-block',
		position: 'relative'
	},
	progress: {
    margin: theme.spacing.unit * 2,
  },
});

const Loading = (props) => {

  const { classes, text } = props;

	return (
		<div className={classes.container}>
		  <CircularProgress className={classes.progress} size={50} />
			{ props.text &&
				<div className={classes.refresh}>{props.text}</div>
			}
		</div>
	)
}

export default withStyles(styles)(Loading);
