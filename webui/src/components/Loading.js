'use strict';

import React from 'react';
import RefreshIndicator from 'material-ui/RefreshIndicator';


const Loading = (props) => {

	const style = {
		container: {
			textAlign: 'center',
			marginTop: 60,
		},
		refresh: {
			display: 'inline-block',
			position: 'relative',
		},
		text: {
			marginTop: 10
		}
	};


	return (
		<div style={ style.container }>
			<RefreshIndicator
				size={40}
				left={10}
				top={0}
				status="loading"
				style={ style.refresh }
			/>
			{ props.text &&
				<div style={ style.text }>{ props.text }</div>
			}
		</div>
	)

}

export default Loading;
