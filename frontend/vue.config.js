const MomentLocalesPlugin = require("moment-locales-webpack-plugin");

module.exports = {
  publicPath: "",
  pages: {
    index: {
      entry: "src/main.ts",
      template: "index.html"
    }
  },
  css: {
    extract: false
  },
  configureWebpack: {
    plugins: [
      // To strip all locales except “en”
      new MomentLocalesPlugin()
    ]
  }
};
