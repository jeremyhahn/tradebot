'use strict';

const express = require('express');
const app = express();
const path = require('path');

app.set( 'port', process.env.PORT || 1234 );

// make the entire contents of public directory accessible
app.use( express.static(
	path.join(__dirname, '../', 'public'),
	{
		// index: false, // don't look for index.html files in sub directories.
		extensions:['html']
	})
);


// for every request made, if the file doesn't exist, return index.html file.
app.get( '/*', (req, res) => {
	res.sendFile( path.join(__dirname, '../', 'public', 'index.html') );
});

app.listen( app.get('port'), function () {
	console.log('Server running at http://localhost:%s', app.get('port'));
});
