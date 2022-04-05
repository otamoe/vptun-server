/* eslint @typescript-eslint/no-require-imports: 0 */

const fs = require('fs')
const path = require('path')
const webpack = require('webpack')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin')
const packageInfo = require('./package.json')
const { VueLoaderPlugin } = require('vue-loader')
const CopyWebpackPlugin = require("copy-webpack-plugin");
const FaviconsWebpackPlugin = require("favicons-webpack-plugin");
const TerserPlugin = require("terser-webpack-plugin");

// fs.readFileSync("./src/css/variable.scss")
const isDev = process.env.NODE_ENV === 'development'

const IsCssExtract = !isDev

if (!isDev) {
  let version = packageInfo.version.split('.');
  version[version.length - 1] = parseInt(version[version.length - 1], 10) + 1;
  packageInfo.version = version.join('.');
  fs.writeFileSync('./package.json', JSON.stringify(packageInfo, null, 2), { flags: 'utf8' })
}


module.exports = {
  mode: process.env.NODE_ENV,
  devtool: isDev ? 'eval-source-map' : 'source-map',

  entry: {
    entry: [
      path.join(__dirname, 'src/entry.ts'),
    ],
  },

  output: {
    path: path.join(__dirname, 'public'),
    publicPath: '/',
    filename: "assets/[name].[contenthash:6].js",
    chunkFilename: 'assets/[name].chunk.[chunkhash:6].js',
  },

  resolve: {
    modules: [
      path.join(__dirname, 'node_modules'),
      path.join(__dirname, 'src'),
    ],
    alias: {
      // "vue": "@vue/runtime-dom",
      "@": path.join(__dirname, 'src'),
      "@/": path.join(__dirname, 'src/'),
      "@/css": path.join(__dirname, 'src/css'),
    },
    extensions: [
      '.js',
      '.jsx',
      '.ts',
      '.tsx',
      '.json',
      '.vue',
      '.scss',
      '.sass',
      '.css'
    ],
  },
  externals: {},
  // cache: false,
  devServer: {
    inline: true,
    hot: true,
    liveReload: true,
    overlay: true,
    disableHostCheck: true,
    watchOptions: {
      poll: true
    },
    port: 8082,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers": "X-Requested-With, content-type, Authorization"
    },
    historyApiFallback: true,

    // before(app) {
    //   app.get('*', function (req, res) {
    // 			res.setHeader("Content-Type", "text/html")
    // 			res.send(html)
    //   })
    // },
  },

  target: 'web',

  module: {
    rules: [{
        test: /\.vue$/,
        use: [{
          loader: 'vue-loader',
          options: {
            compilerOptions: {
              whitespace: 'condense',
            },
            transpileOptions: {
              transforms: {
                dangerousTaggedTemplateString: true
              }
            },
          },
        }],
      },
      {
        test: /\.html?$/,
        use: [{
          loader: 'html-loader',
        }]
      },
      {
        test: /\.(js|jsx)$/,
        use: [{
            loader: 'babel-loader',
          },
          {
            loader: 'eslint-loader',
          },
        ],
      },
      {
        test: /\.tsx?$/,
        use: [{
            loader: 'ts-loader',
            options: {
              appendTsSuffixTo: [/\.(vue|tsx?)$/],
            }
          },
          {
            loader: 'eslint-loader',
          }
        ],
      },
      {
        test: /(\.css)$/,
        use: [{
            loader: 'style-loader',
          },
          {
            loader: "css-loader",
          },
        ],
      },
      {
        test: /(\.scss)$/,
        use: [{
            loader: 'style-loader',
          },
          {
            loader: "css-loader",
            options: {
              importLoaders: 1,
            },
          },
          {
            loader: 'sass-loader',
            options: {
              // additionalData: 
            }
          },
        ],
      },
      {
        test: /(\.sass)$/,
        use: [{
            loader: 'style-loader',
          },
          {
            loader: "css-loader",
            options: {
              importLoaders: 1,
            },
          },
          {
            loader: 'sass-loader',
            options: {
              indentedSyntax: true,
            },
          },
        ],
      },
      {
        test: /\.(gif|jpe?g|png|webp|svg)\??.*$/,
        use: [{
          loader: 'url-loader',
          options: {
            limit: 4096,
            name: 'assets/images/[name].[ext]?v=[hash:8]',
          },
        }, ],
      },
      {
        test: /\.(woff|svg|eot|ttf|woff2|woff)\??.*$/,
        use: [{
          loader: 'url-loader',
          options: {
            limit: 4096,
            name: 'assets/fonts/[name].[ext]?v=[hash:8]',
          },
        }, ],
      },
    ]
  },

  optimization: {
    minimize: !isDev,
    minimizer: [
      new TerserPlugin(),
      new CssMinimizerPlugin(),
    ],
  },



  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.BannerPlugin(`Name: ${packageInfo.name}\nTheme URI: ${packageInfo.homepage}\nVersion: ${packageInfo.version}\nAuthor: ${packageInfo.author}\nLicense: ${packageInfo.license}\nDescription: ${packageInfo.description}`),
    new webpack.DefinePlugin({
      'process.version': JSON.stringify(packageInfo.version),
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: true,
    }),

    new FaviconsWebpackPlugin({
      logo: path.join(__dirname, 'src/public/logo.png'),
      // Enable caching and optionally specify the path to store cached data
      // Note: disabling caching may increase build times considerably
      cache: true,

      publicPath: '/',

      // Prefix path for generated assets
      prefix: '/',

      outputPath: path.join(__dirname, 'public'),


      // Inject html links/metadata (requires html-webpack-plugin).
      // This option accepts arguments of different types:
      //  * boolean
      //    `false`: disables injection
      //    `true`: enables injection if that is not disabled in html-webpack-plugin
      //  * function
      //    any predicate that takes an instance of html-webpack-plugin and returns either
      //    `true` or `false` to control the injection of html metadata for the html files
      //    generated by this instance.
      // inject: true,

      devMode: 'webapp',
      mode: 'webapp',
      favicons: {
        appName: 'VPTun Server',
        appDescription: 'VPTun Server',
        developerName: 'LianYue',
        developerURL: null, // prevent retrieving from the nearest package.json
        background: '#ddd',
        theme_color: '#333',
        icons: {
          coast: false,
          yandex: false
        }
      }
    }),


    new VueLoaderPlugin(),


    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: path.join(__dirname, 'src/index.template.html'),
      minify: {
        collapseWhitespace: true,
        collapseInlineTagWhitespace: true,
        removeComments: true,
        removeRedundantAttributes: true,
        removeScriptTypeAttributes: true,
        removeStyleLinkTypeAttributes: true,
        useShortDoctype: true,
        html5: true,
      },
      scriptLoading: 'defer',
      xhtml: true,
      inlineSource: '.(js|css)$',
    }),


    new CopyWebpackPlugin({
      patterns: [
        { from: "src/public/robots.txt", to: "robots.txt", force: true },
        { from: "src/public/crossdomain.xml", to: "crossdomain.xml", force: true },
      ],
    }),
  ],
};