'use strict';

import React from 'react';
import Header from 'app/components/Header';
import AuthService from 'app/components/AuthService';

class DefaultLayout extends React.Component {

  constructor(props) {
		super(props)
    this.Auth = new AuthService();
	}

	render() {

		return (
			<div>
        {this.Auth.loggedIn() &&
				  <Header title="Tradebot" />
        }
				<div>
					{ this.props.children }
				</div>
			</div>
		)
	}

}

export default DefaultLayout;
