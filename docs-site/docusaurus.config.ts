import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Anexis Server',
  tagline: 'Open-source cloud file storage server',
  favicon: 'img/favicon.ico',

  url: 'https://gurren-software.github.io',
  baseUrl: '/Anexis-Server/',
  organizationName: 'Gurren-Software',
  projectName: 'Anexis-Server',

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
          editUrl: 'https://github.com/Gurren-Software/Anexis-Server/tree/main/docs-site/',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/anexis.png',
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Anexis Server',
      logo: {
        alt: 'Anexis Logo',
        src: 'img/anexis.png',
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
          href: 'https://github.com/Gurren-Software/Anexis-Server',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      logo: {
        alt: 'Anexis Logo',
        src: 'img/anexis.png',
        href: 'https://github.com/Gurren-Software/Anexis-Server',
      },
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
              href: 'https://github.com/Gurren-Software/Anexis-Server',
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
