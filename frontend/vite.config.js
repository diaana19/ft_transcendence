import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [
    react(),
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
          { src: '/logo192.png', sizes: '192x192', type: 'image/png' },
          { src: '/logo512.png', sizes: '512x512', type: 'image/png' },
        ],
      },
     workbox: {
        cleanupOutdatedCaches: true,
        navigateFallbackDenylist: [/^\/api\//],
        runtimeCaching: [
          // NO rule for /api/files/ — unmatched requests pass straight to the
          // network without the SW wrapping the fetch, so streaming 206 video
          // is handled natively by the browser.
          {
            urlPattern: ({ url, request }) =>
              url.pathname.startsWith('/api/') &&
              !url.pathname.startsWith('/api/files/') &&
              !url.pathname.startsWith('/api/ws/') &&
              !url.pathname.startsWith('/api/auth/') &&
              !url.pathname.startsWith('/api/users/me') &&
              !url.pathname.startsWith('/api/notification') &&
              !url.pathname.startsWith('/api/gdpr') &&
              request.method === 'GET',
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api',
              networkTimeoutSeconds: 3,
              expiration: { maxEntries: 200, maxAgeSeconds: 60 * 60 * 24 },
              cacheableResponse: { statuses: [200] },
            },
          },
        ],
      },
    }),
  ],
  build: { outDir: 'build' },
  server: {
    port: 3000,
    proxy: {
      '/api': { target: 'http://localhost:8000', changeOrigin: true, ws: true },
    },
  },
})