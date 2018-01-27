'use strict';

import React from 'react';
import Dialog from 'material-ui/Dialog';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import Loading from 'app/components/Loading';

class ExchangeModal extends React.Component {

	constructor(props) {
		super(props);
		this.state = {
			title: '',
			url: '',
			description: '',
			invalid_url: false,
			processing: false,
		};
		this.updateField = this.updateField.bind(this);
		this.submit = this.submit.bind(this);
	}

	updateField( field, value ) {
		if ( field != 'url' ) {
			this.setState({ [field]: value });
		} else {
			if( /\s/g.test(value) ) {
				this.setState({ [field]: value, invalid_url: true });
			} else {
				this.setState({ [field]: value, invalid_url: false });
			}
		}
	}

	submit() {
		const title = this.state.title;
		const url = this.state.url;
		const description = this.state.description;
		this.setState({ processing: true, title: '', url: '', description: '', invalid_url: false });

		return DB.addSubReddit({ title, url, description })
		.then( res => {			this.setState({ processing: false });
			this.props.update();
			this.props.close();
		});

	}

	render() {

		const actions = [
			<Button
				label="Cancel"
				primary={ true }
				onTouchTap={ this.props.close }
			/>,
			<Button
				label="Submit"
				primary={ true }
				disabled={ ! this.state.title || ! this.state.url }
				onTouchTap={ this.submit }
			/>,
		];

		return (
			<Dialog
				title="Add Exchange"
				actions={ actions }
				modal={ true }
				open={ this.props.open }>

				<p>Please enter your exchange details below.</p>

				{ this.state.processing &&
					<div>
						<Loading />
					</div>
				}

				{ ! this.state.processing &&
				<div>
					<TextField
						floatingLabelText="API Key"
						hintText="Key"
						fullWidth={true}
						defaultValue={ this.state.title }
						onChange={ (event,newValue) => { this.updateField('title', newValue) } }
						autoFocus={true} />
					<TextField
						floatingLabelText="API Secret"
						hintText="Secret"
						fullWidth={true}
						defaultValue={ this.state.description }
						onChange={ (event,newValue) => { this.updateField('description', newValue) } } />
					<TextField
						floatingLabelText="Passphrase"
						hintText="Passphrase"
						fullWidth={true}
						onChange={ (event, newValue) => { this.updateField('url', newValue) } }
						errorText={ this.state.invalid_url && 'Please enter a valid URL.' }
						defaultValue={ this.state.url } />
				</div>
				}

			</Dialog>
		)

	}

}

export default ExchangeModal;
