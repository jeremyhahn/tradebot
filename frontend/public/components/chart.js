var React = require('react')
//var $ = require('jquery'); //Use to load data from the DB


module.exports = React.createClass({
  getInitialState: function() {
    return {
      Config: {
        Username: "odewahn",
        Password: "amoeba21"
      }
    }
  },
  setField: function(e) {
    var s = this.state.Config
    s[e.target.name] = e.target.value
    this.setState({Config: s})
  },
  saveValues: function(bucket, key, e) {
    console.log("The bucket is: ", bucket, key)
    console.log("Sending ", this.state.Config)
    var msg = {
      Bucket: bucket,
      Key: key,
      Value: JSON.stringify(this.state.Config)
    }
    $.ajax({
      type: "POST",
      url: "/db",
      data: JSON.stringify(msg),
      datatype: "JSON"
    })

  },
  render: function() {
    return (
      <div>
        <h1>Hello, World!</h1>
        Username: <input type="text" defaultValue={this.state.Config.Username} name="Username" onChange={this.setField}/>
        <br/>
        Password: <input type="text" defaultValue={this.state.Config.Password} name="Password" onChange={this.setField}/>
        <br/>
        <button onClick={this.saveValues.bind(this, "Config", "CarinaCredentials")}>Do Something</button>
      </div>
    )
  }
})
