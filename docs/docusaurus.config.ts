import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';
// import search from "docusaurus-lunr-search"
// import redirect from "@docusaurus/plugin-client-redirects"

const config: Config = {
  title: 'Gpt4cli Docs',
  tagline: 'An AI coding engine for large, real-world tasks',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://docs.khulnasoft.com',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
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
          routeBasePath: '/', // Serve the docs at the site's root
          editUrl:
            'https://github.com/khulnasoft/gpt4cli/tree/main/docs/',
        },
        blog: false, // Disable the blog
        theme: {
          customCss: './src/css/custom.css',
        },
        
      } satisfies Preset.Options,
    ],
  ],
  themeConfig: {
    // Replace with your project's social card
    image: 'img/gpt4cli-social-preview.png',
    colorMode: {
      defaultMode: "dark",
    },  
    navbar: {
      title: '',
      logo: {
        alt: 'Gpt4cli Logo',
        src: 'img/gpt4cli-logo-light.png',
        srcDark: 'img/gpt4cli-logo-dark.png',
        href: "https://khulnasoft.com",
        height: "2.7rem",
      },
      items: [
        {
          href: 'https://github.com/khulnasoft/gpt4cli',
          label: 'GitHub',
          position: 'right',
        },
        {
          label: 'Discord',
          href: 'https://discord.gg/khulnasoft',
          position: 'right',
        },
        {
          label: 'X',
          href: 'https://x.com/KhulnaSoft',
          position: 'right',
        },
        {
          label: 'YouTube',
          href: 'https://www.youtube.com/@gpt4cli-ny5ry',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',      
      copyright: `Copyright © ${new Date().getFullYear()} KhulnaSoft, Inc.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },

    algolia: {
      // The application ID provided by Algolia
      appId: 'EG57NOYLYX',
      // Public API key: it is safe to commit it
      apiKey: 'a811f8bcdd87a8b3fe7f22a353b968ef',
      indexName: 'gpt4cli',
    }
  } satisfies Preset.ThemeConfig,

  // plugins: [
  //   search,
  //   // [
  //   //   '@docusaurus/plugin-client-redirects',
  //   //   { redirects: [{ from: '/', to: '/intro'}] },
  //   // ],
  // ]
};

export default config;
