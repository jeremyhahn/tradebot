Number.prototype.formatMoney = function() {
 return '$' + this.toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,');
};

Number.prototype.formatCurrency = function(currency) {
  if(currency == "USD") {
    return '$' + this.toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,');
  }
  return this.toFixed(8)
};
