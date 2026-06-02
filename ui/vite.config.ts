import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { writeFileSync } from "node:fs";
import { resolve } from "node:path";

export default defineConfig({
  plugins: [
    react(),
    {
      // emptyOutDir wipes dist (incl. the committed .gitkeep). Go's
      // `//go:embed all:ui/dist` needs >=1 file to vet/build on a fresh
      // checkout, so recreate the placeholder after every build.
      name: "keep-dist-gitkeep",
      closeBundle() {
        writeFileSync(resolve(__dirname, "dist/.gitkeep"), "");
      },
    },
  ],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:3000",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
