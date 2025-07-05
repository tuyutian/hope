import tsconfigPaths from "vite-tsconfig-paths";
import {resolve} from "path";
import {defineConfig} from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from '@tailwindcss/vite'


const pathResolve = (dir: string): string => {
    return resolve(__dirname, ".", dir);
};

// https://vitejs.dev/config/
export default defineConfig({
    esbuild: {
        target: "node20",
    },
    base: "/",
    plugins: [
        tailwindcss(),
        react(),
        tsconfigPaths(),
    ],
    build: {
        // esbuild 打包更快，但是不能去除 console.log，去除 console 使用 terser 模式
        // minify: "terser",
        sourcemap: false,
        /* terserOptions: {
            compress: {
                drop_console: true,
            }
        }, */
        rollupOptions: {
            output: {
                chunkFileNames: "assets/js/[name]-[hash].js", // 引入文件名的名称
                entryFileNames: "assets/js/[name]-[hash].js", // 包的入口文件名称
                assetFileNames: "assets/[ext]/[name]-[hash].[ext]", // 资源文件像 字体，图片等
            },
        },
    },
    resolve: {
        alias: {
            "~": pathResolve("src"),
            "@": pathResolve("src"),
        },
    },
    server: {
        host: "0.0.0.0",
        port: 9527,
        allowedHosts:[
          "s.sunshine-boy.click",
          "api.sunshine-boy.click",
        ]
    },
});
