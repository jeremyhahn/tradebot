'use strict';

import React from 'react';
import Header from 'app/components/Header';

const DefaultLayout = (props) => {

	return (
		<div>
			<Header title="Tradebot" />
			<div className="component--content">
				{ props.children }
			</div>
		</div>
	)
}

export default DefaultLayout;
