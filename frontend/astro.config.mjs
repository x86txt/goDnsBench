import { defineConfig } from 'astro/config';
import tailwind from '@astrojs/tailwind';
import tailwindcss from '@tailwindcss/postcss';

// https://astro.build/config
export default defineConfig({
  integrations: [
    tailwind({
      applyBaseStyles: false,
    })
  ],
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
