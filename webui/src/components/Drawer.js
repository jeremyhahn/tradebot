'use strict';

import React from 'react';

import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';
import { Link } from 'react-router-dom';
import Avatar from 'material-ui/Avatar';
import Subheader from 'material-ui/Subheader';

const LeftDrawer = (props) => {

	const handleClose = () => {
		return props.change(false);
	}

	return (
		<Drawer
			docked={ false }
			width={ 200 }
			open={ props.open }
			onRequestChange={ (status) => props.change(status) }>

			<Subheader>Navigation</Subheader>

			<MenuItem leftIcon={<Avatar src={"images/avatars/128/portfolio.png"} />}
			    onTouchTap={ handleClose } containerElement={<Link to="/portfolio" />} primaryText="Portfolio" />

			<MenuItem leftIcon={<Avatar src={"images/avatars/128/orders.png"} />}
			    onTouchTap={ handleClose } containerElement={<Link to="/orders" />} primaryText="Orders" />

			<MenuItem leftIcon={<Avatar src={"images/avatars/128/exchange.png"} />}
					    onTouchTap={ handleClose } containerElement={<Link to="/exchanges" />} primaryText="Exchanges" />

			<MenuItem leftIcon={<Avatar src={"images/avatars/128/chart.png"} />}
				    onTouchTap={ handleClose } containerElement={<Link to="/chart" />} primaryText="Chart" />

			<MenuItem leftIcon={<Avatar src={"images/avatars/128/settings.png"} />}
			    onTouchTap={ handleClose } containerElement={<Link to="/settings" />} primaryText="Settings" />

		</Drawer>
	)

}

export default LeftDrawer;
