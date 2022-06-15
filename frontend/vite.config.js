import {
  defineConfig,
  splitVendorChunkPlugin
} from 'vite'
import {
  svelte
} from '@sveltejs/vite-plugin-svelte'


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [splitVendorChunkPlugin(), svelte()],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            return id.toString().split('node_modules/')[1].split('/')[0].toString();
          }
        }
      }
    },
    // chunkSizeWarningLimit: 300,
    minify: process.env.NODE_ENV === "production"
  }
})