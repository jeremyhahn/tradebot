'use strict';

import React from 'react';
import Header from 'app/components/Header';

const DefaultLayout = (props) => {

	return (

		<div>

			{props.user.id != null &&
				<Header title="Tradebot" user={props.user}/>
			}

			<div>
				{ props.children }
			</div>

		</div>
	)
}

export default DefaultLayout;
