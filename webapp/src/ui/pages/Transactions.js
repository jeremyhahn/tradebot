import React from 'react';
import classNames from 'classnames';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import Table, {
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TablePagination,
  TableRow,
  TableSortLabel,
} from 'material-ui/Table';
import Typography from 'material-ui/Typography';
import Paper from 'material-ui/Paper';
import Toolbar from 'material-ui/Toolbar';
import Button from 'material-ui/Button';
import Tooltip from 'material-ui/Tooltip';
import FileUpload from 'material-ui-icons/FileUpload';
import FileDownload from 'material-ui-icons/FileDownload';
import FilterListIcon from 'material-ui-icons/FilterList';
import Refresh from 'material-ui-icons/Refresh';
import { lighten } from 'material-ui/styles/colorManipulator';
import Loading from 'app/components/Loading';
import ImportDialog from 'app/components/dialogs/Import';
import withAuth from 'app/components/withAuth';
import AuthService from 'app/components/AuthService';
import {CopyToClipboard} from 'react-copy-to-clipboard';
import Snackbar from 'material-ui/Snackbar';
import Slide from 'material-ui/transitions/Slide';
import Select from 'material-ui/Select';
import { MenuItem } from 'material-ui/Menu';

var downloader = require('js-file-download');

const columnData = [
  { id: 'date', numeric: false, disablePadding: false, label: 'Date' },
  { id: 'type', numeric: false, disablePadding: false, label: 'Type' },
  { id: 'category', numeric: false, disablePadding: false, label: 'Category' },
  { id: 'currency', numeric: false, disablePadding: false, label: 'Currency' },
  { id: 'network', numeric: false, disablePadding: false, label: 'Network' },
  { id: 'quantity', numeric: true, disablePadding: false, label: 'Quantity' },
  { id: 'price', numeric: true, disablePadding: false, label: 'Unit Price' },
  { id: 'fee', numeric: true, disablePadding: false, label: 'Fee' },
  { id: 'total', numeric: true, disablePadding: false, label: 'Total' },
];

class TransactionsHead extends React.Component {

  createSortHandler = property => event => {
    this.props.onRequestSort(event, property);
  };

  render() {
    const { order, orderBy, numSelected, rowCount } = this.props;

    return (
      <TableHead>
        <TableRow>
          {columnData.map(column => {
            return (
              <TableCell
                key={column.id}
                numeric={column.numeric}
                padding={column.disablePadding ? 'none' : 'default'}
                sortDirection={orderBy === column.id ? order : false}>
                <Tooltip
                  title="Sort"
                  placement={column.numeric ? 'bottom-end' : 'bottom-start'}
                  enterDelay={300}>
                  <TableSortLabel
                    active={orderBy === column.id}
                    direction={order}
                    onClick={this.createSortHandler(column.id)}>
                    {column.label}
                  </TableSortLabel>
                </Tooltip>
              </TableCell>
            );
          }, this)}
        </TableRow>
      </TableHead>
    );
  }
}

TransactionsHead.propTypes = {
  numSelected: PropTypes.number.isRequired,
  onRequestSort: PropTypes.func.isRequired,
  order: PropTypes.string.isRequired,
  orderBy: PropTypes.string.isRequired,
  rowCount: PropTypes.number.isRequired,
};

const styles = theme => ({
  root: {
    flex: 1,
    flexGrow: 1,
    paddingLeft: '1%',
    width: '99%',
    marginTop: '68px'
    //marginTop: theme.spacing.unit * 8,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
  table: {
    width: '100%'
  },
  tableWrapper: {
    overflowX: 'auto',
  },
  currencyIcon: {
    paddingLeft: '5px',
    width: '16px',
    height: '16px',
    float: 'right'
  },
  networkIcon: {
    paddingLeft: '5px',
    height: '20px',
    width: '48px'
  },
  loadingContainer: {
    marginTop: '68px'
  },
  leftIcon: {
    marginRight: theme.spacing.unit,
    color: 'grey'
  },
  buttonText: {
    color: 'grey',
    fontSize: '12px'
  }
});

class Transactions extends React.Component {

  constructor(props, context) {
    super(props, context);
    this.Auth = new AuthService();
    this.state = {
      loading: true,
      importDialog: false,
      local_currency: this.Auth.getUser().local_currency,
      order: 'asc',
      orderBy: 'date',
      selected: [],
      data: [],
      page: 0,
      rowsPerPage: 10,
      copied: false,
      copiedValue: '',
      infoMessage: '',
      syncDisabled: false,
      categories: [
        "deposit",
        "withdrawal",
        "trade",
        "income",
        "gift",
        "mining",
        "spend",
        "donation",
        "lost",
        "transfer"
      ],
      selectValues: {}
    }
    this.importDialogHandler = this.importDialogHandler.bind(this)
    this.handleImportDialogClose = this.handleImportDialogClose.bind(this)
    this.appendImportData = this.appendImportData.bind(this)
    this.fetchTransactions = this.fetchTransactions.bind(this)
    this.exportTransactions = this.exportTransactions.bind(this)
    this.syncTransactions = this.syncTransactions.bind(this)
    this.snackbarClose = this.snackbarClose.bind(this)
    this.infoSnackbarClose = this.infoSnackbarClose.bind(this)
  }

  handleRequestSort = (event, property) => {
    const orderBy = property;
    let order = 'desc';

    if (this.state.orderBy === property && this.state.order === 'desc') {
      order = 'asc';
    }

    const data =
      order === 'desc'
        ? this.state.data.sort((a, b) => (b[orderBy] < a[orderBy] ? -1 : 1))
        : this.state.data.sort((a, b) => (a[orderBy] < b[orderBy] ? -1 : 1));

    this.setState({ data, order, orderBy });
  };

  handleChangePage = (event, page) => {
    this.setState({ page });
  };

  handleChangeRowsPerPage = event => {
    this.setState({ rowsPerPage: event.target.value });
  };

  isSelected = id => this.state.selected.indexOf(id) !== -1;

  componentDidMount() {
    this.fetchTransactions()
	}

  fetchTransactions() {
    this.Auth.fetchTransactions()
      .then(function (response) {
        console.log(response);
        if(response.success) {
  		    this.setState({
            loading: false,
            data: response.payload
          })
        }
      }.bind(this))
  }

  exportTransactions() {
    this.Auth.exportTransactions()
      .then(function (response) {
        console.log(response);
        if(response.success) {
          downloader(response.payload, "8949.csv")
        }
      }.bind(this))
  }

  syncTransactions() {
    this.setState({
      syncDisabled: true,
      loading: true
    })
    this.Auth.syncTransactions()
      .then(function (response) {
        console.log(response);
        this.setState({
          syncDisabled: false,
          loading: false
        })
        if(response.success) {
          /*
          var numTxs = response.payload.length
          if(numTxs) {
            this.setState({
              infoSnackbar: true,
              infoMessage: 'Added ' + numTxs + ' new transactions'
            })
            this.fetchTransactions()
          }*/
          this.fetchTransactions()
        }
      }.bind(this))
  }

  currencyIcon(currency) {
    if(currency == undefined) {
      return "images/crypto/128/default.png"
    }
    return "images/crypto/128/" + currency.toLowerCase() + ".png";
  }

  networkIcon(currency) {
    if(currency == undefined) {
      return "images/exchange/default.png"
    }
    return "images/exchange/" + currency.toLowerCase() + ".svg";
  }

  snackbarClose() {
    this.setState({copied: false})
  }

  snackbarTransitionUp(props) {
    return <Slide direction="up" {...props} />;
  }

  infoSnackbarClose() {
    this.setState({infoSnackbar: false})
  }

  appendImportData(newData) {
    console.log(newData)
    /*
    this.setState({
      data: this.state.data.concat(newData),
      order: 'desc'
    })
    this.handleRequestSort(null, 'date')*/
    this.fetchTransactions()
  }

  importDialogHandler() {
    this.setState({importDialog: true})
  }

  handleImportDialogClose() {
    this.setState({importDialog: false})
  }

  handleCategoryChange = (event, value) => {
    var cat = event.target.value
    const formData = new FormData();
    formData.append('category', cat)
    this.Auth.updateCategory(event.target.name, formData)
      .then(function (response) {
        if(response.data.success) {
          for(var i=0; i<this.state.data.length; i++) {
            if(this.state.data[i].id.indexOf(event.target.name) !== -1) {
              this.state.data[i].category = cat.charAt(0).toUpperCase() + cat.slice(1)
              this.setState({data: this.state.data})
            }
          }
        }
      }.bind(this))
  };

  render() {
    const { classes } = this.props;
    const { data, order, orderBy, selected, rowsPerPage, page } = this.state;
    const emptyRows = rowsPerPage - Math.min(rowsPerPage, data.length - page * rowsPerPage);

    return (
      <Paper className={classes.root}>
        {this.state.loading &&
          <div className={classes.loadingContainer}>
            <Loading text="Downloading transactions..." />
          </div>
        }
        <div className={classes.tableWrapper}>
          <Toolbar>
            <Button className={classes.buttonText} size="small" color="inherit" onClick={this.importDialogHandler}>
              Import <FileUpload className={classes.leftIcon} />
            </Button>
            <Button className={classes.buttonText} size="small" color="inherit" onClick={this.exportTransactions}>
              Statement <FileDownload className={classes.leftIcon} />
            </Button>
            <Button disabled={this.state.syncDisabled} className={classes.buttonText} size="small" color="inherit" onClick={this.syncTransactions}>
              API Sync <Refresh className={classes.leftIcon} />
            </Button>
          </Toolbar>
          <Table className={classes.table}>
            <TransactionsHead
              numSelected={selected.length}
              order={order}
              orderBy={orderBy}
              onSelectAllClick={this.handleSelectAllClick}
              onRequestSort={this.handleRequestSort}
              rowCount={data.length}
            />
            <TableBody className={classes.tableBody}>
              {data.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map(n => {
                const isSelected = this.isSelected(n.id);
                return (
                  <TableRow key={"tr" + n.id}>
                    <TableCell padding="none">{new Date(n.date).customFormat()}</TableCell>

                    <TableCell>{n.type}</TableCell>

                    <TableCell>
                      <Select
                        value={n.category.toLowerCase()}
                        onChange={this.handleCategoryChange}
                        inputProps={{
                          id: n.id,
                          name: n.id
                        }}>
                        { this.state.categories.map( category =>
                          <MenuItem key={category} value={category}>{category.charAt(0).toUpperCase() + category.slice(1)}</MenuItem>
                        )}
                      </Select>
                    </TableCell>

                    <TableCell>{n.currency_pair.base}-{n.currency_pair.quote}</TableCell>

                    <TableCell>{n.network_display_name}</TableCell>

                    <TableCell numeric>{n.quantity.formatCurrency(n.quantity_currency)}
                      {(n.fiat_quantity != "0.00" && n.fiat_quantity != "") &&
                      <CopyToClipboard text={n.fiat_quantity} onCopy={()=>this.setState({copied: true, copiedValue: n.fiat_quantity})}>
                        <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.quantity_currency)}
                           title={n.fiat_quantity_currency + " " + n.fiat_quantity.formatCurrency(n.fiat_quantity_currency)} />
                      </CopyToClipboard>
                      ||
                      <img className={classes.currencyIcon} src={this.currencyIcon(n.quantity_currency)} title={n.quantity_currency} />
                      }
                    </TableCell>

                    <TableCell numeric>{n.price.formatCurrency(n.price_currency)}
                      {(n.fiat_price != "0.00" && n.fiat_price != "") &&
                      <CopyToClipboard text={n.fiat_price} onCopy={()=>this.setState({copied: true, copiedValue: n.fiat_price})}>
                        <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.price_currency)}
                           title={n.fiat_price_currency + " " + n.fiat_price.formatCurrency(n.fiat_price_currency) + " per coin"} />
                      </CopyToClipboard>
                      ||
                      <img className={classes.currencyIcon} src={this.currencyIcon(n.price_currency)} title={n.price_currency} />
                      }
                    </TableCell>

                    <TableCell numeric>{n.fee.formatCurrency(n.fee_currency)}
                      {(n.fiat_fee != "0.00" && n.fiat_fee != "") &&
                      <CopyToClipboard text={n.fiat_fee} onCopy={()=>this.setState({copied: true, copiedValue: n.fiat_fee})}>
                        <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.fee_currency)}
                           title={n.fiat_fee_currency + " " + n.fiat_fee.formatCurrency(n.fiat_fee_currency)} />
                      </CopyToClipboard>
                      ||
                      <img className={classes.currencyIcon} src={this.currencyIcon(n.fee_currency)} title={n.fee_currency} />
                      }
                    </TableCell>

                    <TableCell numeric>{n.total.formatCurrency(n.total_currency)}
                      {(n.fiat_total != "0.00" && n.fiat_total != "") &&
                      <CopyToClipboard text={n.fiat_total} onCopy={()=>this.setState({copied: true, copiedValue: n.fiat_total})}>
                        <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.total_currency)}  // fiat_total_currency
                           title={n.fiat_total_currency + " " + n.fiat_total.formatCurrency(n.fiat_total_currency)} />
                      </CopyToClipboard>
                      ||
                      <img className={classes.currencyIcon} src={this.currencyIcon(n.total_currency)} title={n.total_currency} />
                      }
                    </TableCell>

                  </TableRow>
                );
              })}
              {emptyRows > 0 && (
                <TableRow style={{ height: 49 * emptyRows }}>
                  <TableCell colSpan={columnData.length} />
                </TableRow>
              )}
            </TableBody>
            <TableFooter>
              <TableRow>
                <TablePagination
                  colSpan={columnData.length}
                  count={data.length}
                  rowsPerPage={rowsPerPage}
                  page={page}
                  backIconButtonProps={{
                    'aria-label': 'Previous Page',
                  }}
                  nextIconButtonProps={{
                    'aria-label': 'Next Page',
                  }}
                  onChangePage={this.handleChangePage}
                  onChangeRowsPerPage={this.handleChangeRowsPerPage}
                />
              </TableRow>
            </TableFooter>
          </Table>
        </div>
        <ImportDialog open={this.state.importDialog} onClose={this.handleImportDialogClose} addData={this.appendImportData} />
        <Snackbar
          open={this.state.copied}
          onClose={this.snackbarClose}
          autoHideDuration={2000}
          transition={this.snackbarTransitionUp}
          SnackbarContentProps={{
            'aria-describedby': 'message-id',
          }}
          message={<span id="message-id">{this.state.copiedValue} copied to clipboard</span>}/>
        <Snackbar
          open={this.state.infoSnackbar}
          onClose={this.infoSnackbarClose}
          autoHideDuration={2000}
          transition={this.snackbarTransitionUp}
          SnackbarContentProps={{
            'aria-info': 'info-id',
          }}
          message={<span id="info-id">{this.state.infoMessage}</span>}/>
      </Paper>
    );
  }
}

Transactions.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withAuth(withStyles(styles)(Transactions));
