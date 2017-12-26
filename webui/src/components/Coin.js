'use strict';

import React from 'react';
import ActionInfo from 'material-ui/svg-icons/action/info'
import { List, ListItem } from 'material-ui/List';
import { withRouter } from 'react-router-dom';
import IconMenu from 'material-ui/IconMenu';
import MenuItem from 'material-ui/MenuItem';
import IconButton from 'material-ui/IconButton';
import MoreVertIcon from 'material-ui/svg-icons/navigation/more-vert';
import Avatar from 'material-ui/Avatar';


const Coin = (props) => {

	const { data, history } = props;

	const handleTap = () => {
		return history.push(data.url);
	}

	return (
    <ListItem key={data.currency}
				primaryText={ data.currency }
				secondaryText={data.available  + " ($" + data.total +")" }
				onTouchTap={ handleTap }
				leftAvatar={<Avatar src={"images/crypto/128/" + data.currency.toLowerCase() + ".png"} />}
				rightIcon={<IconMenu
			      iconButtonElement={<IconButton><MoreVertIcon /></IconButton>}
			      anchorOrigin={{horizontal: 'left', vertical: 'top'}}
			      targetOrigin={{horizontal: 'left', vertical: 'top'}}>
							<MenuItem primaryText="Auto Trade" />
							<MenuItem primaryText="Buy" />
				      <MenuItem primaryText="Sell" />
				      <MenuItem primaryText="View Chart" />
							<MenuItem primaryText="View Details" />
		    	</IconMenu>}/>
	)
}

export default withRouter( Coin );
