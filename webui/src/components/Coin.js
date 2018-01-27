'use strict';

import React from 'react';
import ActionInfo from 'material-ui-icons/Info'
import Menu, { MenuList, MenuItem } from 'material-ui/Menu';
import List, { ListItem, ListItemSecondaryAction, ListItemText } from 'material-ui/List';
import { withRouter } from 'react-router-dom';
import ListIcon from 'material-ui-icons/List';
import IconButton from 'material-ui/IconButton';
import MoreVertIcon from 'material-ui-icons/MoreVert';
import Avatar from 'material-ui/Avatar';

const Coin = (props) => {

	const { data, history } = props;

	const handleTap = () => {
		return history.push(data.url);
	}

	return (
			<ListItem key={data.currency} button>
				<Avatar src={"images/crypto/128/" + data.currency.toLowerCase() + ".png"} />
				<ListItemText primary={data.currency} secondary={data.available  + " ($" + data.total +")" } />
				<ListItemSecondaryAction>
				<IconButton
          aria-label="More"
          aria-owns={props.anchorEl ? data.currency.toLowerCase() + 'coin-menu' : null}
          aria-haspopup="true"
          onClick={props.menuClickHandler}>
          <MoreVertIcon />
        </IconButton>
				<Menu
				  id={data.currency.toLowerCase() + "coin-menu"}
				  anchorEl={props.anchorEl}
					open={Boolean(props.anchorEl)}
					onClose={props.menuCloseHandler}>
				    <MenuItem onClick={props.autoTradeHandler} disabled={data.currency == 'USD'}>Auto Trade</MenuItem>
						<MenuItem onClick={props.buySellHandler}>Buy / Sell</MenuItem>
						<MenuItem disabled={true}>View Chart</MenuItem>
						<MenuItem>View Details</MenuItem>
        </Menu>
				</ListItemSecondaryAction>
			</ListItem>
	)
}

export default withRouter(Coin);
