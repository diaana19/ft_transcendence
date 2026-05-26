import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [react(), 
    VitePWA({
      registerType: 'autoUpdate',

      manifest: {
        name: 'Transcendence',
        short_name: 'Transcendence',

        start_url: '/',
        display: 'standalone',
          theme_color: '#b89ff3',
          background_color: '#ffffff',

        icons: [
          {
            src: '/logo192.png',
            sizes: '192x192',
            type: 'image/png',
          },

          {
            src: '/logo512.png',
            sizes: '512x512',
            type: 'image/png',
          },
        ],
      },
    }),
  ],
  build: {
    outDir: 'build'
  },
  
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
      }
    }
  }
})
