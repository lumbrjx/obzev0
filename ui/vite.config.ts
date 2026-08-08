import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

export default defineConfig({
  plugins: [react()],
  build: {
    // Output directly into the Go embed directory so `go build` picks it up.
    outDir: resolve(__dirname, '../controller/internal/api/ui/dist'),
    emptyOutDir: true,
  },
})
