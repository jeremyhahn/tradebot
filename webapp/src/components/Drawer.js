'use strict';

import React from 'react';
import Drawer from 'material-ui/Drawer';
import ListSubheader from 'material-ui/List/ListSubheader';
import List, {
  ListItem,
  ListItemAvatar,
  ListItemText,
} from 'material-ui/List';
import { MenuList, MenuItem } from 'material-ui/Menu';
import { Link } from 'react-router-dom';
import Avatar from 'material-ui/Avatar';
import { withStyles } from 'material-ui/styles';

const styles = {
  list: {
    width: 250,
  },
  listFull: {
    width: 'auto',
  },
};

const LeftDrawer = (props) => {

	const handleClose = () => {
		return props.change(false);
	}

	const { classes } = props;

  const navigate = page => {
    props.navigate(page)
  }

	return (

		<Drawer open={ props.open } onClose={ handleClose }>

			<List className={classes.root} subheader={<ListSubheader>Navigation {props.user}</ListSubheader>}>

				<ListItem onTouchTap={ handleClose } onClick={function() {navigate('/portfolio')} } button>
					<ListItemAvatar>
					  <Avatar src={"images/avatars/128/portfolio.png"} />
					</ListItemAvatar>
					<ListItemText primary="Portfolio"/>
				</ListItem>

        <ListItem onTouchTap={ handleClose } onClick={function() {navigate('/transactions')}} button>
					<ListItemAvatar>
					  <Avatar src={"images/avatars/128/transactions.png"} />
					</ListItemAvatar>
					<ListItemText primary="Transactions"/>
				</ListItem>

        <ListItem onTouchTap={ handleClose } onClick={function() {navigate('/orders')}} button>
					<ListItemAvatar>
					  <Avatar src={"images/avatars/128/orders.png"} />
					</ListItemAvatar>
					<ListItemText primary="Orders"/>
				</ListItem>

        <ListItem onTouchTap={ handleClose } onClick={function() {navigate('/exchanges')}} button>
					<ListItemAvatar>
					  <Avatar src={"images/avatars/128/exchange.png"} />
					</ListItemAvatar>
					<ListItemText primary="Exchanges"/>
				</ListItem>

        <ListItem onTouchTap={ handleClose } onClick={function() {navigate('/settings')}} button>
					<ListItemAvatar>
					  <Avatar src={"images/avatars/128/settings.png"} />
					</ListItemAvatar>
					<ListItemText primary="Settings"/>
				</ListItem>

        <ListItem onTouchTap={ handleClose } onClick={function() {navigate('/logout')}} button>
					<ListItemAvatar>
					  <Avatar src={"images/avatars/128/logout.png"} />
					</ListItemAvatar>
					<ListItemText primary="Logout"/>
				</ListItem>

			</List>

		</Drawer>
	)

}

export default withStyles(styles)(LeftDrawer);
