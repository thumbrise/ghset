import {defineConfig} from 'vitepress'

export default defineConfig({
  title: 'ghset — Copy GitHub Repository Settings',
  description: 'ghset init --from owner/repo — copy GitHub repository settings in one command. Settings, security, labels, rulesets. Open-source CLI tool.',
  base: '/ghset/',
  sitemap: {
    hostname: 'https://thumbrise.github.io/ghset/',
  },
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/ghset/favicon.svg' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '96x96', href: '/ghset/favicon-96x96.png' }],
    ['link', { rel: 'apple-touch-icon', sizes: '180x180', href: '/ghset/apple-touch-icon.png' }],
    ['meta', { property: 'og:image', content: 'https://thumbrise.github.io/ghset/og-image.png' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'ghset — Copy GitHub Repository Settings' }],
    ['meta', { property: 'og:description', content: 'ghset init --from owner/repo — copy settings, security, labels, and rulesets from any GitHub repo in one command. Open-source Go CLI.' }],
    ['meta', { property: 'og:url', content: 'https://thumbrise.github.io/ghset/' }],
    ['meta', { name: 'twitter:card', content: 'summary' }],
    ['meta', { name: 'twitter:title', content: 'ghset — Copy GitHub Repository Settings' }],
    ['meta', { name: 'twitter:description', content: 'ghset init --from owner/repo — copy settings, security, labels, rulesets from any repo. One command, zero infrastructure.' }],
    ['meta', { name: 'keywords', content: 'github repository settings cli, declarative github config, github repo settings yaml, copy github repo settings, github settings as code, ghset' }],
  ],

  themeConfig: {
    nav: [
      {text: 'Why ghset?', link: '/why'},
      {text: 'Devlog', link: '/devlog/'},
      {text: 'GitHub', link: 'https://github.com/thumbrise/ghset'},
    ],

    sidebar: {
      '/': [
        {
          text: 'Guide',
          items: [
            {text: 'Why ghset?', link: '/why'},
          ],
        },
        {
          text: 'Devlog',
          items: [
            {text: 'About', link: '/devlog/'},
            {text: '#1 — The Graveyard', link: '/devlog/001-the-graveyard'},
          ],
        },
      ],
    },

    socialLinks: [
      {icon: 'github', link: 'https://github.com/thumbrise/ghset'},
    ],

    editLink: {
      pattern: 'https://github.com/thumbrise/ghset/edit/main/docs/:path',
    },

    footer: {
      message: 'Apache 2.0 · Built in public · Contributions welcome',
    },
  },
})
