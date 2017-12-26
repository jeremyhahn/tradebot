'use strict';

const webpack = require('webpack');
const path = require('path');

const BUILD_DIR = path.resolve(__dirname, '../', 'public/scripts');
const APP_DIR = path.resolve(__dirname, '../', 'src');
const PUBLIC_DIR = path.resolve(__dirname, '../', 'public');

const ExtractTextPlugin = require("extract-text-webpack-plugin");
const HTMLWebpackPlugin = require('html-webpack-plugin');

const InlineManifestWebpackPlugin = require('inline-manifest-webpack-plugin');
const CompressionPlugin = require("compression-webpack-plugin");
const OfflinePlugin = require('offline-plugin');



const VENDOR_LIBS = [
	'react',
	'react-dom',
];


const WebpackConfig = {

	entry: {
		bundle: APP_DIR + '/app.js',
		vendor: VENDOR_LIBS
	},
	output: {
		path: BUILD_DIR,
		filename: '[name].[chunkhash].js',
		chunkFilename: '[name].[chunkhash].chunk.js',
		publicPath: '/scripts/',
	},

	module: {
		rules: [
			{
				enforce: 'pre',
				test: /\.js$/,
				exclude: /node_modules/,
				loader: 'eslint-loader',
				include : APP_DIR
			},
			{
				loader: 'babel-loader',
				test: /\.js$/,
				exclude: /node_modules/,
				include : APP_DIR,
				options: {
					presets: [
						[
							'env',
							{ modules: false }
						],
						'react'
					],
					plugins: [ 'syntax-dynamic-import' ]
				}
			},
			{
				use: ExtractTextPlugin.extract({
					use: 'css-loader',
				}),
				test: /\.css$/
			},
			{
				loader: 'json-loader',
				test: /\.json$/
			}
		],
	},

	plugins: [
		new ExtractTextPlugin({
			filename: 'styles.css',
			allChunks: true
		}),
		new webpack.optimize.CommonsChunkPlugin({
			names: ['vendor', 'manifest'],
			minChunks: Infinity,
		}),
		new HTMLWebpackPlugin({
			inject: false,
			filename: '../index.html',
			template: 'scripts/index.ejs',
			minify: {
				collapseBooleanAttributes: true,
				removeComments: true,
				collapseWhitespace: true,
			}
		}),

		new InlineManifestWebpackPlugin({
			name: 'webpackManifest'
		}),
		new webpack.optimize.AggressiveMergingPlugin(),
		new webpack.optimize.UglifyJsPlugin({
			beautify: false,
			mangle: {
				screw_ie8: true,
				// keep_fnames: true
			},
			compress: {
				warnings: false,
				screw_ie8: true
			},
			comments: false
		}),
		new CompressionPlugin({
			asset: "[path].gz[query]",
			algorithm: "gzip",
			test: /\.(js|html|css)$/,
			threshold: 10240,
			minRatio: 0.8
		}),
		new OfflinePlugin({
			caches: 'all',
			relativePaths: false,
			publicPath: '/scripts/',
			ServiceWorker: {
				events: true,
				entry: './scripts/sw-handler.js',
				output: '../sw.js',
				publicPath: '/sw.js',
				navigateFallbackURL: '/'
			},
			AppCache: {
				events: true,
				output: '../appcache',
				publicPath: '/appcache/'
			},
			externals: [
				'/',
				'https://fonts.googleapis.com/css?family=Roboto:300,400,500',
			],
		}),


	],

	resolve: {
		alias: {
			app: APP_DIR,
			public: PUBLIC_DIR
		},
	},


};

module.exports = WebpackConfig;

