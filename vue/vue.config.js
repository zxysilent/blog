const path = require("path");
const TerserPlugin = require("terser-webpack-plugin");
const resolve = dir => {
	return path.join(__dirname, dir);
};
module.exports = {
	// 部署应用时的根路径(默认'/'),也可用相对路径
	publicPath: process.env.NODE_ENV === "development" ? "/" : "/dist/", //修改为你的URL+/dist/",//若不单独部署则不用修改
	// 运行时生成的生产环境构建文件的目录(默认''dist''，构建之前会被清除)
	outputDir: "../dist",
	// 放置生成的静态资源(s、css、img、fonts)的(相对于 outputDir 的)目录(默认'')
	assetsDir: "",
	// 指定生成的 index.html 的输出路径(相对于 outputDir)也可以是一个绝对路径。
	indexPath: "index.html",
	// 生成的静态资源在它们的文件名中包含了 hash 以便更好的控制缓存
	filenameHashing: true,
	// pages 里配置的路径和文件名在你的文档目录必须存在 否则启动服务会报错
	// https://cli.vuejs.org/zh/config/#pages
	pages: {
		index: {
			// 除了 entry 之外都是可选的
			entry: "src/main.js", // page 的入口,每个"page"应该有一个对应的 JavaScript 入口文件
			template: "public/index.html", // 模板来源
			filename: "index.html", // 在 dist/index.html 的输出
			title: "blog", // 当使用 title 选项时,在 template 中使用：<title><%= htmlWebpackPlugin.options.title %></title>
			chunks: ["chunk-vendors", "chunk-common", "index"] // 在这个页面中包含的块，默认情况下会包含,提取出来的通用 chunk 和 vendor chunk
		}
	},
	// 是否在保存的时候检查
	lintOnSave: false,
	// 是否使用包含运行时编译器的 Vue 构建版本。设置为 true 后你就可以在 Vue 组件中使用 template 选项了，但是这会让应用额外增加 10kb 左右。
	runtimeCompiler: false,
	// 生产环境是否生成 sourceMap 文件
	productionSourceMap: false,
	// https://cli.vuejs.org/zh/config/#devserver
	devServer: {
		// 环境配置
		host: "127.0.0.1",
		port: 8080,
		https: false,
		// 没有页面刷新的情况下启用热模块替换
		hotOnly: false,
		open: true //配置自动启动浏览器
		//一个代理
		//proxy: 'http://127.0.0.1:3000'
		// 多个代理
		// proxy: {
		//   "/api": {
		//     target: "http://127.0.0.1:1010",
		//     changeOrigin: true
		//   },
		//   "/foo": {
		//     target: "http://127.0.0.1:1009"
		//   }
		// }
	},
	// 通过 webpack-merge 合并到最终的配置中
	// https://cli.vuejs.org/zh/config/#configurewebpack
	configureWebpack: {
		// https://webpack.js.org/configuration/performance/
		performance: {
			hints: false // 取消打包文件过大的警告
		},
		plugins: [
			// new webpack.ProvidePlugin({
			// 	"window.Quill": "quill"
			// })
		],
		optimization: {
			//https://webpack.js.org/configuration/optimization/#optimization-minimize
			minimizer: [
				new TerserPlugin({
					cache: true,
					terserOptions: {
						compress: {
							warnings: false,
							drop_console: true,
							drop_debugger: true
						}
					}
				})
			]
		}
	},
	// 允许对内部的 webpack 配置进行更细粒度的修改
	// https://cli.vuejs.org/zh/config/#chainwebpack
	chainWebpack: config => {
		config.resolve.alias.set("@", resolve("src"));
	}
};
