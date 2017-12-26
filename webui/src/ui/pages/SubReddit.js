'use strict';

import React from 'react';
import Paper from 'material-ui/Paper';
import _ from 'lodash';

import { List } from 'material-ui/List';
import Subheader from 'material-ui/Subheader';
import SubRedditPost from 'app/components/SubRedditPost';

import DB from 'app/utils/DB';
import Loading from 'app/components/Loading';


class SubReddit extends React.Component {

	constructor(props) {
		super(props);
		const current_subreddit_url = _.replace( this.props.location.pathname, '/', '' );
		this.state = {
			subreddit: {},
			url: current_subreddit_url,
			loading: true,
			error: false,
			posts: [],
		};
		this.findSubRedditByURL = this.findSubRedditByURL.bind(this);
		this.fetchPosts = this.fetchPosts.bind(this);
		this.updatePosts = this.updatePosts.bind(this);
	}


	componentDidMount() {
		this.findSubRedditByURL();
		this.fetchPosts();
	}


	findSubRedditByURL() {
		DB.findSubRedditByURL( this.state.url )
		.then( (record) => {
			if ( record ) {
				this.setState({ subreddit: record });
			}
		});
	}



	updatePosts( raw ) {
		const data = raw.data.children;
		const records = [];
		data.map( record => {
			const data = record.data;
			const rec = {
				id: data.id,
				title: data.title,
				url: data.url,
				content: data.selftext,
				author: data.author,
				ups: data.ups,
				downs: data.downs,
			}
			records.push( rec );
		});
		this.setState({ loading: false, posts: records });
	}



	// https://developers.google.com/web/fundamentals/instant-and-offline/offline-cookbook/
	fetchPosts() {
		let networkDataReceived = false;
		const URL = `https://www.reddit.com/r/${ this.state.url }.json`;
		let this_component = this;

		const networkUpdate = fetch( URL )
			.then( (response) => {
				return response.json();
			})
			.then( (data) => {
				networkDataReceived = true;
				this_component.updatePosts( data );
			});

		// fetch cached data
		caches.match( URL )
			.then(function(response) {
				if ( ! response ) throw Error("No data");
				return response.json();
			})
			.then(function(body) {
				// don't overwrite newer network data
				if ( ! networkDataReceived ) {
					this_component.updatePosts( body );
				}
			})
			.catch(function() {
				// we didn't get cached data, the network is our last hope:
				return networkUpdate;
			})
			.catch(function(error) {
				this_component.setState({ loading: false, error: true });
			});


	}


	render() {

		return (
			<Paper style={{ padding: 20 }} zDepth={1} rounded={false}>
				<h1 style={{ paddingLeft: 16 }}>{ this.state.subreddit.title }</h1>
				{ this.state.subreddit.description && <div style={{ paddingLeft: 16, marginTop: 10 }}>{ this.state.subreddit.description }</div> }

				{ this.state.loading &&
					<Loading />
				}

				{ this.state.error &&
					<div style={{ padding: 20, backgroundColor: '#FFCDD2', marginTop: 25 }}>
						<h3>Error Occoured.</h3>
						<p>Please enter a valid URL.</p>
					</div>
				}


				{ ! this.state.loading && ! this.state.error &&
				<List>
					<Subheader style={{ textTransform: 'uppercase' }}>All Posts</Subheader>
					{ this.state.posts.map( item => <SubRedditPost key={ item.id } data={ item } /> ) }
				</List>
				}


			</Paper>
		)

	}

}

export default SubReddit;

