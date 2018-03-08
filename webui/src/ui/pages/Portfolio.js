'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import List, { ListItem, ListItemText } from 'material-ui/List';
import ListSubheader from 'material-ui/List/ListSubheader';
import Loading from 'app/components/Loading';
import EmptyExchangeList from 'app/components/EmptyExchangeList';
import ExchangeModal from 'app/components/modal/Exchange';
import AutoTradeModal from 'app/components/modal/AutoTrade'
import BuySellDialog from 'app/components/dialogs/BuySell'
import Coin from 'app/components/Coin';
import Avatar from 'material-ui/Avatar';
import Typography from 'material-ui/Typography';
import withAuth from 'app/components/withAuth';

const currencies = [
  {
    value: 'USD',
    label: '$',
  },
  {
    value: 'EUR',
    label: '€',
  },
  {
    value: 'BTC',
    label: '฿',
  },
  {
    value: 'JPY',
    label: '¥',
  },
];

class Portfolio extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			anchorEl: null,
			addExchangeModal: false,
			autoTradeModal: false,
			sellModal: false,
			loading: true,
			portfolio: {
				exchanges: []
			},
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
		this.loadPortfolio = this.loadPortfolio.bind(this);
		this.handleMenuClick = this.handleMenuClick.bind(this);
		this.handleMenuClose = this.handleMenuClose.bind(this);
	}

	componentDidMount() {
		this.loadPortfolio();
	}

	componentWillUnmount() {
	 this.ws.close();
  }

	loadPortfolio() {
		var loc = window.location, new_uri;
		var protocol = (loc.protocol === "https:") ? "wss" : "ws";
		if(this.ws == null) {
		  this.ws = new WebSocket(protocol + "://localhost:8080/ws/portfolio");
	  }
		var _this = this;
		var ws = this.ws;
		this.ws.onopen = function() {
      ws.send(JSON.stringify(_this.props.user))
		};
		this.ws.onclose = function() {
			console.log("Websocket connection closed")
		}
		this.ws.onmessage = evt => {
			 var portfolio = JSON.parse(evt.data);
			 console.log(portfolio);

			 // Server-side JSON marshalling turns empty array into null
			 if(portfolio.exchanges == null) {
				 portfolio.exchanges = [];
			 }
			 if(portfolio.wallets == null) {
				 portfolio.wallets = [];
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

	handleMenuClick = event => {
		console.log('handleMenuClick')
		this.setState({ anchorEl: event.currentTarget })
	}

	handleMenuClose = () => {
		console.log('handleMenuClose')
		this.setState({ anchorEl: null })
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

			<div style={{flex: '100%', height: '100%'}}>

				<div style={{float: 'right', marginTop: '20px', marginRight: '35px'}}>
					<Typography type="subheading" gutterBottom>Net worth: { this.state.netWorth.formatMoney() }</Typography>
			  </div>

				<Paper style={{ marginTop: '60px', height: '100%'}}>

					{ this.state.portfolio.exchanges.map( exchange =>

						<List key={exchange.name}>
							<ListSubheader style={{ textTransform: 'uppercase' }}>
								{ exchange.name + " - " + exchange.satoshis + " BTC - " + exchange.total.formatMoney()}
							</ListSubheader>

							{
								exchange.coins.map( coin =>
									<Coin buySellHandler={this.handleBuySellModalOpen}
												autoTradeHandler={this.handleAutoTradeModalOpen}
												key={ coin.currency } data={ coin }
												anchorEl={this.state.anchorEl}
												menuClickHandler={this.handleMenuClick}
												menuCloseHandler={this.handleMenuClose} />
							)}
					  </List>
					)}

					{ this.state.portfolio.wallets &&

						<List>
						<ListSubheader style={{ textTransform: 'uppercase' }}>Wallets</ListSubheader>
						{ this.state.portfolio.wallets.map( (wallet, i) =>
							<ListItem key={wallet.currency + "-" + i} button>
								<Avatar src={"images/crypto/128/" + wallet.currency.toLowerCase() + ".png"} />
								<ListItemText primary={wallet.currency} secondary={wallet.balance  + " (" + wallet.value.formatMoney() +")" } />
							</ListItem>
						)}
						</List>

					}

					<BuySellDialog
						open={ this.state.sellModal }
						close={ this.handleBuySellModalClose }
						update={ this.handleBuySellModalUpdate } />

					<ExchangeModal
						open={ this.state.addExchangeModal }
						close={ this.handleAddExchangeModalClose }
						update={ this.handleAddExchangeModalUpdate } />

					<AutoTradeModal
						open={ this.state.autoTradeModal }
						close={ this.handleAutoTradeModalClose }
						update={ this.handleAutoTradeModalUpdate } />

				</Paper>

			</div>
		)
	}

}

export default withAuth(Portfolio);
