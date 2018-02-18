'use strict';

import React from 'react';
import { withStyles } from 'material-ui/styles';
import { CircularProgress } from 'material-ui/Progress';

const styles = theme => ({
	container: {
		textAlign: 'center',
		marginTop: '50%',
		marginLeft: '50%'
	},
	refresh: {
		display: 'inline-block',
		position: 'relative',
	},
	text: {
		marginTop: '50%'
	},
	progress: {
		margin: `0 ${theme.spacing.unit * 2}px`,
	}
});

const Loading = (props) => {

	return (
		<div style={ styles.container }>
		  <CircularProgress className={styles.progress} size={50} />
			{ props.text &&
				<div style={ styles.text }>{ props.text }</div>
			}
		</div>
	)

}

export default withStyles(styles)(Loading);
