import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Anexis Server',
  tagline: 'Open-source cloud file storage server',
  favicon: 'img/favicon.ico',

  url: 'https://anexis-server.github.io',
  baseUrl: '/',
  organizationName: 'Treefle-labs',
  projectName: 'anexis-server',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/Treefle-labs/anexis-server/tree/main/docs-site/',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/anexis-social-card.png',
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Anexis Server',
      logo: {
        alt: 'Anexis Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'docSidebar',
          sidebarId: 'apiSidebar',
          position: 'left',
          label: 'API',
        },
        {
          type: 'docSidebar',
          sidebarId: 'deploymentSidebar',
          position: 'left',
          label: 'Deployment',
        },
        {
          href: 'https://github.com/Treefle-labs/anexis-server',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Documentation',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/getting-started',
            },
            {
              label: 'Configuration',
              to: '/docs/configuration',
            },
            {
              label: 'Deployment',
              to: '/docs/deployment/docker',
            },
          ],
        },
        {
          title: 'API',
          items: [
            {
              label: 'Authentication',
              to: '/docs/api/auth',
            },
            {
              label: 'Files',
              to: '/docs/api/files',
            },
            {
              label: 'Links',
              to: '/docs/api/links',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/Treefle-labs/anexis-server',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Anexis Server. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;