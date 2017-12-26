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
import DB from 'app/utils/DB';
import Loading from 'app/components/Loading';

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
//		this.loadExchanges = this.loadExchanges.bind(this);
	}

	componentDidMount() {
		this.loadExchanges();
	}

	loadExchanges() {
		/*
		DB.findAllSubReddits()
		.then( (res) => {
			if ( res ) {
				this.setState({ loading: false, exchanges: res });
				console.log(res)
			} else {
				this.setState({ loading: false });
			}
		});
		*/
		var ws = new WebSocket("ws://localhost:8080/portfolio");
		ws.onopen = function() {
			 ws.send(JSON.stringify({currency: "BTC-USD"}));
		};
		ws.onmessage = evt => {
			 var exchangeList = JSON.parse(evt.data);
			 console.log(exchangeList);
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
			<Paper style={{ padding: 20, }} zDepth={1} rounded={false}>

        { this.state.exchanges.map( exchange =>
					<List key={exchange.name}>
					  <Subheader style={{ textTransform: 'uppercase' }}>{exchange.name + " - $" + exchange.total}</Subheader>
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
		)
	}

}

export default Portfolio;
