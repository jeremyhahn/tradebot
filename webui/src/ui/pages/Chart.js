
import React from 'react';
import { render } from 'react-dom';
import TradeChart from 'app/components/TradeChart';
import { getData } from "./utils"
import Paper from 'material-ui/Paper';

import { TypeChooser } from "react-stockcharts/lib/helper";

class Chart extends React.Component {

  constructor(props) {
		super(props);
		this.setState({
			data: []
		})
	}

	componentDidMount() {
		getData().then(data => {
			this.setState({ data })
		})
	}

	render() {
		return (
			<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>
				<h2>Chart</h2>
				<TypeChooser>
				{type => <TradeChart type={type} data={this.state.data} />}
				</TypeChooser>
			</Paper>
		)
	}
}

export default Chart;
