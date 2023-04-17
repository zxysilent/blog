import { resolve } from "path";
import { defineConfig, splitVendorChunkPlugin } from "vite";
import vue from "@vitejs/plugin-vue";
/* 路径方法 */
const pathResolver = (pathStr: string): string => {
    return resolve(__dirname, ".", pathStr);
};

const endpoint = (mode: string): string => {
    if (mode === "development") {
        return "http://127.0.0.1:8085";
    }
    if (mode === "devonline") {
        return "https://blog.zxysilent.com";
    }
    return "";
}

// https://vitejs.dev/config/
export default ({ mode }) => {
    return defineConfig({
        base: mode === "production" ? "/dist/" : "/",
        plugins: [vue(), splitVendorChunkPlugin()],
        resolve: {
            // alias: {
            //     '@': pathResolver('./src'),
            // },
            alias: [
                {
                    find: "@",
                    replacement: pathResolver("src") + "/",
                },
            ],
        },
        server: {
            open: true,
            port: 8082,
            proxy: {
                // /api=>xxx/api
                // '/api': endpoint(mode),
                // '/static': endpoint(mode),
                '/api': {
                    target: endpoint(mode),
                    changeOrigin: true
                },
                '/static': {
                    target: endpoint(mode),
                    changeOrigin: true
                },
            },
        },
        define: {
            __APP_INFO__: JSON.stringify({}),
        },
        esbuild: {
            //@ts-ignore
            minify: true,
            treeShaking: true,
            minifySyntax: true,
            minifyWhitespace: true,
            minifyIdentifiers: true,
            drop: ['console', 'debugger'],

        },
        build: {
            outDir: "./dist",
            assetsDir: "static",
            minify: "esbuild",
            emptyOutDir: true,
            chunkSizeWarningLimit: 1792,

        },
        css: {
            preprocessorOptions: {
                less: {
                    modifyVars: {},
                    javascriptEnabled: true,
                    additionalData: `@import "src/styles/var.less";`,
                },
            },
        },
    });
};
