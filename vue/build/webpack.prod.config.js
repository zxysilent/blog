const webpack = require("webpack");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const cleanWebpackPlugin = require("clean-webpack-plugin");
const merge = require("webpack-merge");
const webpackBaseConfig = require("./webpack.base.config.js");
const fs = require("fs");
const path = require("path");
const package = require("../package.json");

fs.open("./build/env.js", "w", function(err, fd) {
  const buf = 'export default "production";';
  //fs.write(fd, buf, 0, buf.length, 0, function(err, written, buffer) {});
  fs.write(fd, buf, 0, 'utf-8', function(err, written, buffer) {});
});

module.exports = merge(webpackBaseConfig, {
  output: {
    publicPath: "http://127.0.0.1:88/res/dist/dist/", //修改为你的URL+/dist/
    filename: "[name].[hash].js",
    chunkFilename: "[name].[hash].chunk.js"
  },
  plugins: [
    new cleanWebpackPlugin(["dist/*"], {
      root: path.resolve(__dirname, "../../res/") //default '../'
    }),
    new ExtractTextPlugin({
      filename: "[name].[hash].css",
      allChunks: true
    }),
    new webpack.optimize.CommonsChunkPlugin({
      // name: 'vendors',
      // filename: 'vendors.[hash].js'
      name: ["vender-exten", "vender-base"],
      minChunks: Infinity
    }),
    new webpack.DefinePlugin({
      "process.env": {
        NODE_ENV: '"production"'
      }
    }),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: true,
        drop_console: true,
        drop_debugger: true
      }
    }),
    new CopyWebpackPlugin([
      {
        from: "favicon.ico"
      }
    ]),
    new HtmlWebpackPlugin({
      title: "vue-blog v" + package.version,
      favicon: "./favicon.ico",
      filename: "../index.html",
      template: "!!ejs-loader!./src/template/index.ejs",
      inject: false
    })
  ]
});
