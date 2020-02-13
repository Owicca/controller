const path = require('path'),
  webpack = require('webpack'),
  outputPath = 'public/',
  outputJsFileName = 'main.js',
  jsDir = path.resolve(__dirname, 'js/');


module.exports = {
  mode: 'development',
  watch: true,
  watchOptions: {
    aggregateTimeout: 300,
    poll: 1000,
    ignored: /node_modules/
  },
  entry:  path.resolve(jsDir, 'main.js'),
  output: {
    path: path.resolve(__dirname, outputPath),
    filename: outputJsFileName,
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader']
      }
    ]
  },
  resolve: {
    extensions: ['*', '.js', '.jsx']
  },
  plugins: [
    new webpack.ProgressPlugin(),
  ]
};
