'use strict';

import React from 'react';

import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';
import { Link } from 'react-router-dom';


const LeftDrawer = (props) => {

	const handleClose = () => {
		return props.change(false);
	}

	return (
		<Drawer
			docked={ false }
			width={ 200 }
			open={ props.open }
			onRequestChange={ (status) => props.change(status) }
		>
			<MenuItem onTouchTap={ handleClose } containerElement={<Link to="/portfolio" />} primaryText="Portfolio" />
			<MenuItem onTouchTap={ handleClose } containerElement={<Link to="/settings" />} primaryText="Settings" />
		</Drawer>
	)

}

export default LeftDrawer;
