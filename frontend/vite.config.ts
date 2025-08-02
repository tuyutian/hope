import tsconfigPaths from "vite-tsconfig-paths";
import { resolve } from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const pathResolve = (dir: string): string => {
  return resolve(__dirname, ".", dir);
};

// https://vitejs.dev/config/
export default defineConfig({
  esbuild: {
    target: "node20",
  },
  base: "/",
  plugins: [tailwindcss(), react(), tsconfigPaths()],
  build: {
    sourcemap: false,
    rollupOptions: {
      output: {
        chunkFileNames: "assets/js/[name]-[hash].js",
        entryFileNames: "assets/js/[name]-[hash].js",
        assetFileNames: "assets/[ext]/[name]-[hash].[ext]",
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
    allowedHosts: ["s.protectifyapp.com", "api.protectifyapp.com"],
  },
});
