import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    'intro',
    'getting-started',
    'installation',
    'configuration',
    'architecture',
  ],
  apiSidebar: [
    'api/auth',
    'api/files',
    'api/links',
    'api/migration',
    'api/backup',
  ],
  deploymentSidebar: [
    'deployment/docker',
    'deployment/docker-compose',
    'deployment/production',
    'deployment/self-hosted',
  ],
};

export default sidebars;