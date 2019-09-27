const path = require('path');

function root(file) {
  return path.resolve(__dirname, file);
}

module.exports = {
  entry: {
    index: './index.js',
  },

  output: {
    path: root('build'),
    filename: '[name].bundle.js'
  },

  resolve: {
    extensions: ['.js'],
    modules: ['node_modules']
  },

  module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            cacheDirectory: true,
            presets: ['@babel/preset-env'],
            plugins: [
              '@babel/plugin-proposal-class-properties',
              ["@babel/plugin-transform-runtime",
                {
                  "regenerator": true
                }
              ]
            ]
          },
        }
      },
      {
        test: /\.css?$/,
        use: ['to-string-loader', 'css-loader'],

      },
      {
        test: /\.html$/,
        use: [{
          loader: 'html-loader',
          options: {
            minimize: true
          }
        }]
      },
      {
        test: /\.(png|jpe?g|gif|svg|woff|woff2|ttf|eot|ico)(\?.*)?$/,
        loader: 'file-loader?name=./assets/[name].[hash].[ext]'
      }
    ]
  }
}
