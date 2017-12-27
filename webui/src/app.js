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
import DefaultLayout from 'app/ui/layouts/Default';
import Portfolio from 'app/ui/pages/Portfolio';
import Exchanges from 'app/ui/pages/Exchanges';
import Settings from 'app/ui/pages/Settings';
import Chart from 'app/ui/pages/Chart';
import { install } from 'offline-plugin/runtime';

render(
	(
	<Router>
	<MuiThemeProvider>

		<DefaultLayout>
			<Route exact path="/" component={ Portfolio } />
			<Switch>
				<Route exact path="/portfolio" component={ Portfolio } />
				<Route exact path="/exchanges" component={ Exchanges } />
				<Route exact path="/chart" component={ Chart } />
				<Route exact path="/settings" component={ Settings } />
			</Switch>
		</DefaultLayout>

	</MuiThemeProvider>
	</Router>
	),
	document.getElementById('root')
);

// install the service worker.
install();
