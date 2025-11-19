 /*global process */

import { fileURLToPath, URL } from 'url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
   define: {
      // enable hydration mismatch details in production build
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true'
   },
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
   css: {
      preprocessorOptions : {
          scss: {
              api: "modern-compiler",
              additionalData: `@use "@/assets/theme/colors.scss" as *;`
          },
      }
   },
})


