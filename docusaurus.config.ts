import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Tìm Hiểu Thánh Kinh',
  tagline: 'Lời Chúa là ngọn đèn cho chân tôi, Ánh sáng cho đường-lối tôi.',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://phankhoi.vercel.app',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'facebook', // Usually your GitHub org/user name.
  projectName: 'VI1934', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'vi',
    locales: ['vi'],
  },
  
  presets: [
    [
      'classic',  {
        docs: {
          sidebarPath: './sidebars.ts',
          path: "VI1934",
          routeBasePath: "/",
          sidebarCollapsible: true,
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/docusaurus-social-card.jpg',
    navbar: {
      title: 'Tìm Hiểu Thánh Kinh',
      logo: {
        alt: 'My Site Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'ĐỌC',
        },
        {
          href: 'https://phucam.tv/',
          label: 'PhucAm.tv',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Tài Liệu',
          items: [
            {
              label: 'Tín Lý',
              to: 'https://tinly.phucam.tv/',
            },
            {
              label: 'Lẽ thật Ngày Sa-bát',
              to: 'https://sabat.phucam.tv/',
            },
          ],
        },
        {
          title: 'Cộng Đồng',
          items: [
            {
              label: 'Facebook',
              href: 'https://www.facebook.com/phucamtv',
            },
            {
              label: 'Youtube',
              href: 'https://www.youtube.com/@70lan7',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/70lan7',
            },
          ],
        },
        {
          title: 'Nội mạng',
          items: [
            {
              label: 'PhucAm.tv',
              href: 'https://phucam.tv',
            },
          ],
        },
      ],
      copyright: `Tìm Hiểu Thánh Kinh`,
    },
    prism: {
      theme: prismThemes.jettwaveLight,
      darkTheme: prismThemes.oneDark,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
