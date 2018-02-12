'use strict';

import React from 'react';
import withAuth from 'app/components/withAuth'
import Header from 'app/components/Header';
/*
const DefaultLayout = (props) => {

	return (
		<div>
		  {props.user != null &&
				<Header title="Tradebot" />
			}
			<div>
				{ props.children }
			</div>
		</div>
	)
}

export default DefaultLayout;
*/

class DefaultLayout extends React.Component {

  constructor(props) {
		super(props)
	}

	render() {

		return (
			<div>
				<Header title="Tradebot" />
				<div>
					{ this.props.children }
				</div>
			</div>
		)
	}

}

export default DefaultLayout;
