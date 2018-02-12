'use strict';

import React from 'react';
import { render } from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'

import injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();

import 'app/css/reset.css';
import 'app/css/style.css';
import 'app/css/helper.css';
import 'app/css/typography.css';
import 'app/util/array.js'
import 'app/util/number.js'

import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import createMuiTheme from 'material-ui/styles/createMuiTheme'
import createPalette from 'material-ui/styles/createPalette'
import {grey, amber, red} from 'material-ui/colors'

import DefaultLayout from 'app/ui/layouts/Default';
import Portfolio from 'app/ui/pages/Portfolio';
import Exchanges from 'app/ui/pages/Exchanges';
import Settings from 'app/ui/pages/Settings';
import Chart from 'app/ui/pages/Chart';
import OrderHistory from 'app/ui/pages/OrderHistory';

import Login from 'app/components/Login';
import Register from 'app/ui/pages/Register';
import AuthService from 'app/components/AuthService';
import withAuth from 'app/components/withAuth';

import { install } from 'offline-plugin/runtime';

const Auth = new AuthService();

const muiTheme = createMuiTheme({
  palette: {
    primary: {
      light: '#7986cb',
      main: '#2196F3',
      dark: '#303f9f',
      contrastText: 'rgba(255, 255, 255, 1)',
    },
    secondary: {
      light: '#ff4081',
      main: '#f50057',
      dark: '#c51162',
      contrastText: 'rgba(255, 255, 255, 1)',
    },
  },
});

const handleLogout = function() {
   Auth.logout()
   this.props.history.replace('/login');
}

const user = {
  id: null,
  username: null,
  local_currency: null
}

render(
	(
	<Router>
	<MuiThemeProvider theme={muiTheme}>
		<DefaultLayout user={user}>
			<Route exact path="/" component={ user.id == null ? Login : Portfolio } />
			<Switch>
				<Route exact path="/portfolio" component={ Portfolio } />
				<Route exact path="/orders" component={ OrderHistory } />
				<Route exact path="/exchanges" component={ Exchanges } />
				<Route exact path="/chart" component={ Chart } />
				<Route exact path="/settings" component={ Settings } />
        <Route exact path="/login" component={ Login } />
        <Route exact path="/register" component={ Register } />
			</Switch>
		</DefaultLayout>
	</MuiThemeProvider>
	</Router>
	),
	document.getElementById('root')
);

install();
