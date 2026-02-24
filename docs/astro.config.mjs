import { defineConfig } from 'astro/config';
import config from "./config.mjs"
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: config.url,
  base: config.basePath,
  integrations: [
    starlight({
      title: config.title,
      description: config.description,
      social: {
        github: config.github,
      },
      editLink: {
        baseUrl: config.githubDocs,
      },

      sidebar: [
        {
          label: 'Getting Started',
          items: [
            { label: 'Introduction', link: '/' },
            { label: 'Installation', link: '/guides/installation' },
            { label: 'Quick Start', link: '/guides/quickstart' },
            { label: 'Configuration', link: '/guides/configuration' },
          ],
        },
        {
          label: 'API Reference',
          items: [
            { label: 'app', link: '/api/app' },
            { label: 'config', link: '/api/config' },
            { label: 'domain', link: '/api/domain' },
            { label: 'errors', link: '/api/errors' },
            { label: 'package', link: '/api/package' },
            { label: 'ui', link: '/api/ui' },
            { label: 'utils', link: '/api/utils' },
            { label: 'workflow', link: '/api/workflow' },
            {
              label: 'adapters',
              collapsed: true,
              items: [
                { label: 'clipboard', link: '/api/adapters/clipboard' },
                { label: 'editor', link: '/api/adapters/editor' },
              ],
            },
          ],
        },
      ],
      customCss: ['./src/styles/custom.css'],
    }),
  ],
});
