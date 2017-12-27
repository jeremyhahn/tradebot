import React from 'react';
import {Tabs, Tab} from 'material-ui/Tabs';
import FlatButton from 'material-ui/FlatButton';
import Menu from 'material-ui/Menu';
import MenuItem from 'material-ui/MenuItem';
import Dialog from 'material-ui/Dialog';
import Loading from 'app/components/Loading';

const styles = {
  headline: {
    fontSize: 24,
    paddingTop: 16,
    marginBottom: 12,
    fontWeight: 400,
  }
};

class BuySellModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      slideIndex: 0
    };
  }

  render() {

		const actions = [
			<FlatButton
				label="Cancel"
				primary={ true }
				onTouchTap={ this.props.close }
			/>,
			<FlatButton
				label="Submit"
				primary={ true }
				disabled={ ! this.state.title || ! this.state.url }
				onTouchTap={ this.submit }
			/>,
		];

    return (

      <Dialog
				title="Sell"
				actions={ actions }
				modal={ true }
				open={ this.props.open }>

				{ this.state.processing &&
					<div>
						<Loading />
					</div>
				}

				{ ! this.state.processing &&
          <Tabs>
    		    <Tab label="Market" >
    		      <div>
    		        <h2 style={styles.headline}>Tab One</h2>
    		        <p>
    		          Create a new market order.
    		        </p>
    		      </div>
    		    </Tab>
    		    <Tab label="Limit" >
    		      <div>
    		        <h2 style={styles.headline}>Tab Two</h2>
    		        <p>
    		          TODO: Create new limit order.
    		        </p>
    		      </div>
    		    </Tab>
    		    <Tab label="Stop" data-route="/home">
    		      <div>
    		        <h2 style={styles.headline}>Tab Three</h2>
    		        <p>
    		          TODO: Create a new stop order.
    		        </p>
    		      </div>
    		    </Tab>
    		  </Tabs>
			  }

			</Dialog>
    );
  }
}

export default BuySellModal;
