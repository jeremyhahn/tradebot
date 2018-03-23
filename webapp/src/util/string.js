String.prototype.formatCurrency = function(currency) {
  if(currency == "USD") {
    return '$' + this
    //return '$' + this.formatFiat()
  } /*else {
    return this.toString().match(/^-?\d+(?:\.\d{0,8})?/)[0];
  }*/
  return this
};
