'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import { List } from 'material-ui/List';
import Subheader from 'material-ui/Subheader';
import Coin from 'app/components/Coin';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import ContentAdd from 'material-ui/svg-icons/content/add';
import Avatar from 'material-ui/Avatar'
import EmptyExchangeList from 'app/components/EmptyExchangeList';
import ExchangeModal from 'app/components/modal/Exchange';
import Loading from 'app/components/Loading';

const netWorthDiv = {
	float: 'right',
	marginRight: '35px'
}

const netWorthHeader = {
	fontWeight: 'bold'
}

class Portfolio extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			modal: false,
			loading: true,
			exchanges: [],
		};
		this.handleModalOpen = this.handleModalOpen.bind(this);
		this.handleModalClose = this.handleModalClose.bind(this);
		this.handleModalUpdate = this.handleModalUpdate.bind(this);
		this.loadExchanges = this.loadExchanges.bind(this);
	}

	componentDidMount() {
		this.loadExchanges();
	}

	componentWillUnmount() {
	 this.ws.close();
  }

	loadExchanges() {
		if(this.ws == null) {
		  this.ws = new WebSocket("ws://localhost:8080/portfolio");
	  }
		var ws = this.ws
		this.ws.onopen = function() {
			 ws.send(JSON.stringify({currency: "BTC-USD"}));
		};
		this.ws.onmessage = evt => {
			 var exchangeList = JSON.parse(evt.data);
			 //console.log(exchangeList);
			 this.setState({
				 loading: false,
				 exchanges: exchangeList
			 })
		};

	}

	handleModalOpen() {
		this.setState({ modal: true });
	}

	handleModalClose() {
		this.setState({ modal: false });
	}

	handleModalUpdate() {
		this.loadExchanges();
	}

	render() {

		if ( this.state.loading ) {
			return <Loading text="Loading coins..." />
		}

		if ( this.state.exchanges.length < 1 ) {
			return (
				<div>

					<EmptyExchangeList openModal={ this.handleModalOpen } />

					<ExchangeModal
						open={ this.state.modal }
						close={ this.handleModalClose }
						update={ this.handleModalUpdate } />

				</div>
			)
		}

		return (

			<div>

				<div style={netWorthDiv}>
			    <Subheader style={netWorthHeader}>Net worth: { this.state.exchanges.sum("total").formatMoney() }</Subheader>
			  </div>


				<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>

	        { this.state.exchanges.map( exchange =>

						<List key={exchange.name}>
						  <Subheader style={{ textTransform: 'uppercase' }}>{
								exchange.name + " - " + exchange.satoshis + " BTC - " + exchange.total.formatMoney() }
							</Subheader>
						  { exchange.coins.map( coin => <Coin key={ coin.currency } data={ coin } /> ) }
					  </List>
					)}

					<FloatingActionButton style={{ position: 'fixed', bottom: 50, right: 50 }} onTouchTap={ this.handleModalOpen }>
						<ContentAdd />
					</FloatingActionButton>

					<ExchangeModal
						open={ this.state.modal }
						close={ this.handleModalClose }
						update={ this.handleModalUpdate } />

				</Paper>

			</div>
		)
	}

}

export default Portfolio;
