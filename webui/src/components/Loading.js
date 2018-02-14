'use strict';

import React from 'react';
import { CircularProgress } from 'material-ui/Progress';

const Loading = (props) => {

	const style = theme => ({
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

	return (
		<div style={ style.container }>
		  <CircularProgress className={style.progress} size={50} />
			{ props.text &&
				<div style={ style.text }>{ props.text }</div>
			}
		</div>
	)

}

export default Loading;
