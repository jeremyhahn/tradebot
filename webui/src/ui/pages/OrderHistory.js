'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import Loading from 'app/components/Loading';
import { TablePagination } from 'react-pagination-table';
import axios from 'axios';

class OrderHistory extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			orders: [],
      page: 0,
      rowsPerPage: 25,
		};
	}

	componentDidMount() {
		var self = this
		axios.get('/orderhistory')
      .then(function (response) {
        console.log(response);
		    self.setState({ orders: response.data })
      })
	}

	componentWillUnmount() {
  }

	render() {

    const Header = ["Exchange", "Date", "Type", "Currency", "Quantity", "Price"];
    const { orders, page, rowsPerPage } = this.state;
    //const emptyRows = rowsPerPage - Math.min(rowsPerPage, this.state.orders.length - page * rowsPerPage);

		if ( this.state.loading ) {
			return <Loading text="Loading orders..." />
		}

		if ( this.state.orders.length < 1 ) {
			return (
					<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>
 						No orders yet, try placing a new limit order.
					</Paper>
			)
		}

		return (

      <Paper>
      <div>
        <TablePagination
          title="Order History"
          headers={Header}
          data={orders}
          columns="exchange.date.type.currency.quantity.price"
          perPageItemCount={rowsPerPage}
          totalCount={ orders.length }
          arrayOption={ [["size", 'all', ' ']] }
        />
      </div>
      </Paper>
		)
	}

}

export default OrderHistory;
