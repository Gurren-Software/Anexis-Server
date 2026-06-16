/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx,md,mdx}',
    './docs/**/*.{md,mdx}',
    './docusaurus.config.{js,ts}',
  ],
  darkMode: ['selector', '[data-theme="dark"]'],
  corePlugins: {
    preflight: false,
  },
  theme: {
    extend: {
      colors: {
        anexis: {
          ink: 'var(--anexis-ink)',
          muted: 'var(--anexis-muted)',
          page: 'var(--anexis-page)',
          surface: 'var(--anexis-surface)',
          elevated: 'var(--anexis-elevated)',
          border: 'var(--anexis-border)',
          primary: 'var(--anexis-primary)',
          primarySoft: 'var(--anexis-primary-soft)',
          primaryStrong: 'var(--anexis-primary-strong)',
          code: 'var(--anexis-code)',
        },
      },
      boxShadow: {
        anexis: 'var(--anexis-shadow)',
        glow: 'var(--anexis-glow)',
      },
      backgroundImage: {
        'anexis-hero': 'var(--anexis-hero-gradient)',
        'anexis-band': 'var(--anexis-band-gradient)',
        'anexis-text': 'var(--anexis-text-gradient)',
      },
    },
  },
  plugins: [],
};
