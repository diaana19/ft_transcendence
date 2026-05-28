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

      workbox: {
        cleanupOutdatedCaches: true,
        runtimeCaching: [
          // Files (avatars, post media): each has a UUID id and never changes,
          // so cache-first with a long expiry is safe and gives the biggest
          // offline win — once an avatar is loaded, it works without network.
          {
            urlPattern: ({ url }) => url.pathname.startsWith('/api/files/'),
            handler: 'CacheFirst',
            options: {
              cacheName: 'files',
              expiration: {
                maxEntries: 500,
                maxAgeSeconds: 60 * 60 * 24, // 1 day
              },
              cacheableResponse: { statuses: [0, 200] },
            },
          },
          // Other public GETs under /api/. Network-first with a 3 s timeout:
          // fresh data when online, last-known cached response when offline.
          // User-specific endpoints are excluded — Workbox's default cache key
          // is (url, method) and does NOT include Authorization, so caching
          // /api/users/me etc. would let one user see another user's data on
          // the same device after a logout. Clear "api" cache on logout if
          // you ever broaden this rule.
          {
            urlPattern: ({ url, request }) =>
              url.pathname.startsWith('/api/') &&
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
              expiration: {
                maxEntries: 200,
                maxAgeSeconds: 60 * 60 * 24, // 1 day
              },
              cacheableResponse: { statuses: [200] },
            },
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
