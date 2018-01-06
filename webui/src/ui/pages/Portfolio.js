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
import AutoTradeModal from 'app/components/modal/AutoTrade'
import BuySellModal from 'app/components/modal/BuySell'

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
			addExchangeModal: false,
			autoTradeModal: false,
			sellModal: false,
			loading: true,
			portfolio: {},
			netWorth: 0.0
		};
		this.handleAddExchangeModalOpen = this.handleAddExchangeModalOpen.bind(this);
		this.handleAddExchangeModalClose = this.handleAddExchangeModalClose.bind(this);
		this.handleAddExchangeModalUpdate = this.handleAddExchangeModalUpdate.bind(this);
		this.handleAutoTradeModalOpen = this.handleAutoTradeModalOpen.bind(this);
		this.handleAutoTradeModalClose = this.handleAutoTradeModalClose.bind(this);
		this.handleAutoTradeModalUpdate = this.handleAutoTradeModalUpdate.bind(this);
		this.handleBuySellModalOpen = this.handleBuySellModalOpen.bind(this);
		this.handleBuySellModalClose = this.handleBuySellModalClose.bind(this);
		this.handleBuySellModalUpdate = this.handleBuySellModalUpdate.bind(this);
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
			 ws.send(JSON.stringify({user: {id: 1, username: "jhahn"}}));
		};
		this.ws.onmessage = evt => {
			 var portfolio = JSON.parse(evt.data);
			 console.log(portfolio);

			 // Server-side JSON marshalling turns empty array into null
			 if(portfolio.exchanges == null) {
				 portfolio.exchanges = [];
			 }
			 for(var i=0; i<portfolio.exchanges.length; i++) {
				 if(portfolio.exchanges[i].coins == null) {
					 portfolio.exchanges[i].coins = [];
				 }
			 }

			 this.setState({
				 loading: false,
				 portfolio: portfolio,
				 netWorth: portfolio.net_worth
			 })
		};
	}

	handleAddExchangeModalOpen() {
		this.setState({ addExchangeModal: true });
	}
	handleAddExchangeModalClose() {
		this.setState({ addExchangeModal: false });
	}
	handleAddExchangeModalUpdate() {
		console.log("handleAddExchangeModalUpdate() fired")
	}

	handleAutoTradeModalOpen(e) {
		e.preventDefault();
		this.setState({ autoTradeModal: true });
	}
	handleAutoTradeModalClose() {
		this.setState({ autoTradeModal: false });
	}
	handleAutoTradeModalUpdate() {
		console.log("handleAutoTradeModalUpdate() fired")
	}

	handleBuySellModalOpen(e) {
		e.preventDefault()
		this.setState({ sellModal: true });
	}
	handleBuySellModalClose() {
		this.setState({ sellModal: false });
	}
	handleBuySellModalUpdate() {
		console.log("handleBuySellModalModalUpdate() fired")
	}

	render() {

		if ( this.state.loading ) {
			return <Loading text="Loading coins..." />
		}

		if ( this.state.portfolio.exchanges.length < 1 ) {
			return (
				<div>

					<EmptyExchangeList openModal={ this.handleAddExchangeModalOpen } />

					<ExchangeModal
						open={ this.state.addExchangeModal }
						close={ this.handleAddExchangeModalClose }
						update={ this.handleAddExchangeModalUpdate } />

				</div>
			)
		}

		return (

			<div>

				<div style={netWorthDiv}>
			    <Subheader style={netWorthHeader}>Net worth: { this.state.netWorth.formatMoney() }</Subheader>
			  </div>

				<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>

	        { this.state.portfolio.exchanges.map( exchange =>

						<List key={exchange.name}>
						  <Subheader style={{ textTransform: 'uppercase' }}>{
								exchange.name + " - " + exchange.satoshis + " BTC - " + exchange.total.formatMoney() }
							</Subheader>
						  {
									exchange.coins.map( coin =>
									  <Coin buySellHandler={this.handleBuySellModalOpen}
										      autoTradeHandler={this.handleAutoTradeModalOpen}
													key={ coin.currency } data={ coin } /> )
							}
					  </List>
					)}

					{ this.state.portfolio.wallets.map( wallet =>

						<List key={wallet.currency}>
						  <Subheader style={{ textTransform: 'uppercase' }}>{
								wallet.currency + " Wallet - " + wallet.balance + " - " + wallet.net_worth.formatMoney() }
							</Subheader>
					  </List>
					)}

					<ExchangeModal
						open={ this.state.addExchangeModal }
						close={ this.handleAddExchangeModalClose }
						update={ this.handleAddExchangeModalUpdate } />

					<AutoTradeModal
						open={ this.state.autoTradeModal }
						close={ this.handleAutoTradeModalClose }
						update={ this.handleAutoTradeModalUpdate } />

					<BuySellModal
						open={ this.state.sellModal }
						close={ this.handleBuySellModalClose }
						update={ this.handleBuySellModalUpdate } />

				</Paper>

			</div>
		)
	}

}

export default Portfolio;
