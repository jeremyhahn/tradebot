import React from 'react';
import {Tabs, Tab} from 'material-ui/Tabs';
import Dialog from 'material-ui/Dialog';
import Button from 'material-ui/Button';
import TextField from 'material-ui/TextField';
import { MenuList, MenuItem } from 'material-ui/Menu';
import Select from 'material-ui/Select';
import Checkbox from 'material-ui/Checkbox';
import Loading from 'app/components/Loading';

const styles = {
  headline: {
    fontSize: 24,
    paddingTop: 16,
    marginBottom: 12,
    fontWeight: 400,
  },
	indicatorsLabel: {
		paddingBotton: 25
  },
  checkbox: {
    marginBottom: 16,
  }
};

const indicators = [
  'RSI',
  'MACD',
  'Bollinger Bands',
  'EMA'
];

class AutoTradeModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      slideIndex: 0,
			open: false,
			rsiChecked: false,
			bollingerChecked: false,
			macdChecked: false,
			emaChecked: false,
			values: []
    };
  }

	menuItems(values) {
    return indicators.map((name) => (
      <MenuItem
        key={name}
        insetChildren={true}
        checked={values && values.indexOf(name) > -1}
        value={name}
        primaryText={name}
      />
    ));
  }

	handleSelectMenuChange(event, index, values) {
		this.setState({values})
	}

  handleChange(event, index, value) {
	   this.setState({value});
	}

	updateRsiCheck() {
    this.setState((oldState) => {
      return {
        rsiChecked: !oldState.checked,
      };
    });
  }

	updateBollingerCheck() {
    this.setState((oldState) => {
      return {
        bollingerChecked: !oldState.checked,
      };
    });
  }

	updateMacdCheck() {
    this.setState((oldState) => {
      return {
        macdChecked: !oldState.checked,
      };
    });
  }

  render() {

		 const {values} = this.state;

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
				title="Auto Trade"
				actions={ actions }
				modal={ true }
				open={ this.props.open }>

				{ this.state.processing &&
					<div>
						<Loading />
					</div>
				}

				{ ! this.state.processing &&
				<div>
					<Select floatingLabelText="Trading Strategy">
						<MenuItem value={1} primaryText="Position Trading" />
						<MenuItem value={2} primaryText="Swing Trading" />
					</Select>
					<div>
						<p style={styles.indicatorsLabel}>Which financial indicators would you like to use?</p>
					  <Select multiple={true} hintText="Financial Indicator" value={values} onChange={this.handleSelectMenuChange} >
              {this.menuItems(values)}
						</Select>
					</div>
				</div>
			  }

			</Dialog>
    );
  }
}

export default AutoTradeModal;
