/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        // Tech Innovation Theme
        'electric-blue': '#0066ff',
        'neon-cyan': '#00ffff',
        'dark-gray': '#1e1e1e',
        'light-gray': '#2a2a2a',
        'lighter-gray': '#3a3a3a',
      },
      fontFamily: {
        sans: ['DejaVu Sans', 'system-ui', 'sans-serif'],
        mono: ['DejaVu Sans Mono', 'monospace'],
      },
    },
  },
  plugins: [],
}
