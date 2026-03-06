import { defineConfig } from 'astro/config';
import config from "./config.mjs"
import sidebar from "./sidebar.mjs"
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
      customCss: ['./src/styles/custom.css'],
      expressiveCode: {
        themes: ['material-theme-lighter', 'material-theme-darker'],
      },
      sidebar: sidebar,
    }),
  ],
});
