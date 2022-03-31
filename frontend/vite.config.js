 /*global process */

import { fileURLToPath, URL } from 'url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
   plugins: [vue()],
   resolve: {
      alias: {
         '@': fileURLToPath(new URL('./src', import.meta.url))
      }
   },
   server: { // this is used in dev mode only
      port: 8080,
      proxy: {
        '/api': {
          target: process.env.COLLECT_SRV, // export COLLECT_SRV=http://localhost:8085
          changeOrigin: true,
          logLevel: 'debug'
        },
        '/version': {
          target: process.env.COLLECT_SRV,
          changeOrigin: true,
          logLevel: 'debug'
        },
        '/healthcheck': {
          target: process.env.COLLECT_SRV,
          changeOrigin: true,
          logLevel: 'debug'
        },
      }
   },
})


