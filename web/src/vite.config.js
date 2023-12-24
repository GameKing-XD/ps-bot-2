import { defineConfig } from "vite";

export default defineConfig({
        base: "/app",
        server: {
                proxy: {
                        '/api': 'http://localhost:8080'
                }
        },
        build: {
                outDir: "../assets/",
                emptyOutDir: true,
        }
})
