'use strict';

import React from 'react';
import { render } from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

// import all the custom styles
import 'app/css/reset.css';
import 'app/css/style.css';
import 'app/css/helper.css';
import 'app/css/typography.css';
import 'app/util/array.js'
import 'app/util/number.js'

import injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();

import DefaultLayout from 'app/ui/layouts/Default';
import Portfolio from 'app/ui/pages/Portfolio';
import Settings from 'app/ui/pages/Settings';
import ChartPage from 'app/ui/pages/ChartPage';
import { install } from 'offline-plugin/runtime';

render(
	(
	<Router>
	<MuiThemeProvider>

		<DefaultLayout>
			<Route exact path="/" component={ Portfolio } />
			<Switch>
				<Route exact path="/portfolio" component={ Portfolio } />
				<Route exact path="/settings" component={ Settings } />
				<Route exact path="/chart" component={ ChartPage } />
			</Switch>
		</DefaultLayout>

	</MuiThemeProvider>
	</Router>
	),
	document.getElementById('root')
);

// install the service worker.
install();
