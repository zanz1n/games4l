import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { imagetools } from "vite-imagetools";
import { VitePWA } from "vite-plugin-pwa";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        svelte(),
        imagetools(),
        VitePWA({
            manifest: {
                name: "Memories for Life",
                short_name: "Memories4Life",
                description: "Um jogo desenvolvido pelo Colégio Poliedro em parceria com o Hospital Reger",
                lang: "pt-BR",
                theme_color: "#ffffff",
                icons: [
                    {
                        src: "favicon.png",
                        sizes: "128x128",
                        type: "image/png"
                    },
                    {
                        src: "favicon-big.png",
                        sizes: "512x512",
                        type: "image/png"
                    },
                    {
                        src: "favicon-big.png",
                        sizes: "512x512",
                        type: "image/png",
                        purpose: "any"
                    },
                    {
                        src: "favicon-big.png",
                        sizes: "512x512",
                        type: "image/png",
                        purpose: "maskable"
                    }
                ]
            },
            registerType: "autoUpdate",
            injectRegister: "inline",
            workbox: {
                globPatterns: ["**/*.{js,css,html,ico,png,svg,gif,webp,avif,json,webm}"],
                inlineWorkboxRuntime: true,
            }
        })
    ],
    build: {
        target: "esnext"
    }
});
