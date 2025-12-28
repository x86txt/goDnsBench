import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/postcss';

// https://astro.build/config
export default defineConfig({
  integrations: [],
  output: 'static',
  outDir: './dist',
  build: {
    assets: 'assets'
  },
  vite: {
    css: {
      postcss: {
        plugins: [tailwindcss()],
      },
    },
    build: {
      target: 'esnext',
      minify: 'esbuild',
      cssMinify: true
    }
  }
});
