'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import { List } from 'material-ui/List';
import Subheader from 'material-ui/Subheader';
import Avatar from 'material-ui/Avatar'
import Loading from 'app/components/Loading';
import {
  Table,
  TableBody,
  TableHeader,
  TableHeaderColumn,
  TableRow,
  TableRowColumn,
} from 'material-ui/Table';
import axios from 'axios';

class OrderHistory extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			orders: []
		};
	}

	componentDidMount() {
		var self = this
		axios.get('/orderhistory')
		//axios.get('http://localhost:8080/orderhistory')
      .then(function (response) {
        console.log(response);
		    self.setState({ orders: response.data })
      })
	}

	componentWillUnmount() {
  }

	render() {

		if ( this.state.loading ) {
			return <Loading text="Loading orders..." />
		}

		if ( this.state.orders.length < 1 ) {
			return (
					<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>
 						No orders yet, try placing a new stop or limit order.
					</Paper>
			)
		}

		return (

			<div>

				<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>
				<Table>
 					<TableHeader>
	 					<TableRow>
		 					<TableHeaderColumn>Date</TableHeaderColumn>
							<TableHeaderColumn>Exchange</TableHeaderColumn>
		 					<TableHeaderColumn>Currency</TableHeaderColumn>
		 					<TableHeaderColumn>Quantity</TableHeaderColumn>
							<TableHeaderColumn>Price</TableHeaderColumn>
	 					</TableRow>
 					</TableHeader>
 				  <TableBody>
	        { this.state.orders.length > 0 && this.state.orders.map( order =>
						<TableRow>
			        <TableRowColumn>{order.date}</TableRowColumn>
							<TableRowColumn>{order.exchange}</TableRowColumn>
				      <TableRowColumn>{order.currency_pair.base} - {order.currency_pair.quote}</TableRowColumn>
							<TableRowColumn>{order.quantity}</TableRowColumn>
							<TableRowColumn>{order.price}</TableRowColumn>
			      </TableRow>
					)}
					</TableBody>
  			</Table>
				</Paper>

			</div>
		)
	}

}

export default OrderHistory;
