Date.prototype.customFormat = function() {
  var hours = this.getHours();
  var minutes = this.getMinutes();
  var ampm = hours >= 12 ? 'pm' : 'am';
  hours = hours % 12;
  hours = hours ? hours : 12;
  minutes = minutes < 10 ? '0' + minutes : minutes;
  return (this.getMonth()+1) + "-" + this.getDate() + "-" + this.getFullYear() + " " +
    hours + ':' + minutes + ':' + this.getSeconds() + ' ' + ampm;
}
