Number.prototype.formatMoney = function() {
 return '$' + this.toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,');
};

Number.prototype.addMoneySymbol = function() {
  var decimals = this.toString().split('.')
  var len = (decimals.length > 1) ? decimals[1].length : 0
 return '$' + this.toFixed(len).replace(/(\d)(?=(\d{3})+\.)/g, '$1,');
};
