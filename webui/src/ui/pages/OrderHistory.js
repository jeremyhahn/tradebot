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
import IconButton from 'material-ui/IconButton';
import Tooltip from 'material-ui/Tooltip';
import DeleteIcon from 'material-ui-icons/Delete';
import FilterListIcon from 'material-ui-icons/FilterList';
import { lighten } from 'material-ui/styles/colorManipulator';
import withAuth from 'app/components/withAuth';
import AuthService from 'app/components/AuthService';

const columnData = [
  { id: 'date', numeric: false, disablePadding: true, label: 'Date' },
  { id: 'exchange', numeric: false, disablePadding: false, label: 'Exchange' },
  { id: 'type', numeric: false, disablePadding: false, label: 'Type' },
  { id: 'currency_pair', numeric: true, disablePadding: false, label: 'Currency' },
  { id: 'quantity', numeric: true, disablePadding: false, label: 'Quantity' },
  { id: 'price', numeric: true, disablePadding: false, label: 'Price' },
  { id: 'fee', numeric: true, disablePadding: false, label: 'Fee' },
  { id: 'total', numeric: true, disablePadding: false, label: 'Total' }
];

class OrderHistoryHead extends React.Component {

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

OrderHistoryHead.propTypes = {
  numSelected: PropTypes.number.isRequired,
  onRequestSort: PropTypes.func.isRequired,
  order: PropTypes.string.isRequired,
  orderBy: PropTypes.string.isRequired,
  rowCount: PropTypes.number.isRequired,
};

const styles = theme => ({
  root: {
    flex: 1,
    paddingLeft: '1%',
    width: '99%',
    marginTop: '68px'
    //marginTop: theme.spacing.unit * 8,
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
  }
});

class OrderHistory extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.Auth = new AuthService();
    this.state = {
      local_currency: this.Auth.getUser().local_currency,
      order: 'asc',
      orderBy: 'date',
      selected: [],
      data: [],
      page: 0,
      rowsPerPage: 10,
    };
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
    this.Auth.fetch('/api/v1/orderhistory')
      .then(function (response) {
        console.log(response);
        if(response.success) {
          for(var i=0; i<response.payload.length; i++) {
            response.payload[i].price = response.payload[i].price;
          }
  		    this.setState({ data: response.payload })
        }
      }.bind(this))
	}

  currencyIcon(currency) {
    return "images/crypto/128/" + currency.toLowerCase() + ".png";
  }

  render() {
    const { classes } = this.props;
    const { data, order, orderBy, selected, rowsPerPage, page } = this.state;
    const emptyRows = rowsPerPage - Math.min(rowsPerPage, data.length - page * rowsPerPage);

    return (
      <Paper className={classes.root}>
        <div className={classes.tableWrapper}>
          <Table className={classes.table}>
            <OrderHistoryHead
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
                  <TableRow key={n.id}>
                    <TableCell padding="none">{n.date}</TableCell>
                    <TableCell numeric>{n.exchange}</TableCell>
                    <TableCell numeric>{n.type}</TableCell>
                    <TableCell numeric>{n.currency_pair.base}-{n.currency_pair.quote}</TableCell>
                    <TableCell numeric>{n.quantity}</TableCell>
                    <TableCell numeric>{n.price.formatCurrency(n.currency_pair.quote == this.state.local_currency ? n.currency_pair.quote : n.currency_pair.base)}
                      <img className={classes.currencyIcon}
                         src={this.currencyIcon(n.currency_pair.quote == this.state.local_currency ? n.currency_pair.quote : n.currency_pair.base)}
                         title={n.currency_pair.quote == this.state.local_currency ? n.currency_pair.quote : n.currency_pair.base} />
                    </TableCell>
                    <TableCell numeric>{n.fee.formatCurrency(n.currency_pair.quote)}
                      <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.currency_pair.base)}
                           title={n.currency_pair.base} />
                    </TableCell>
                    <TableCell numeric>{n.total.formatCurrency(n.currency_pair.quote)}
                      <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.currency_pair.quote)}
                           title={n.currency_pair.quote} />
                    </TableCell>
                  </TableRow>
                );
              })}
              {emptyRows > 0 && (
                <TableRow style={{ height: 49 * emptyRows }}>
                  <TableCell colSpan={6} />
                </TableRow>
              )}
            </TableBody>
            <TableFooter>
              <TableRow>

                <div><i className="material-icons">file_download</i>Download 8949 Statement</div>

                <TablePagination
                  colSpan={6}
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

      </Paper>
    );
  }
}

OrderHistory.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withAuth(withStyles(styles)(OrderHistory));
