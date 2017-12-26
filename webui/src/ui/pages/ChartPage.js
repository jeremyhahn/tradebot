
import React from 'react';
import { render } from 'react-dom';
import TradeChart from 'app/components/TradeChart';
import { getData } from "./utils"
import Paper from 'material-ui/Paper';

import { TypeChooser } from "react-stockcharts/lib/helper";

class ChartComponent extends React.Component {
	componentDidMount() {
    /*
		getData().then(data => {
			this.setState({ data })
		})*/
	}
	render() {
		if (this.state == null) {
			return <div>Loading...</div>
		}
		return (
			<ChartComponent />
		)
	}
}
