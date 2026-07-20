import { defineConfig, type DefaultTheme } from 'vitepress'
import { cpSync, existsSync, readdirSync, readFileSync, statSync } from 'node:fs'
import { basename, dirname, join, relative } from 'node:path'

const repo = 'https://github.com/justCDQ/Lili'
const root = process.cwd()

const sections = [
  { text: '前端深化', dir: '01-frontend' },
  { text: '产品能力', dir: '02-product' },
  { text: '交互设计', dir: '03-interaction-design' },
  { text: 'AI 工程', dir: '04-ai' },
  { text: '后端与数据', dir: '05-backend-data' }
]

function titleFromMarkdown(file: string) {
  const content = readFileSync(file, 'utf8')
  const match = content.match(/^#\s+(.+)$/m)
  return match?.[1]?.trim() || basename(file, '.md')
}

function linkFromFile(file: string) {
  const path = relative(root, file).replaceAll('\\', '/')
  return `/${path.replace(/\.md$/, '')}`
}

function directoryItems(dir: string): DefaultTheme.SidebarItem[] {
  const fullDir = join(root, dir)

  return readdirSync(fullDir)
    .filter((entry) => !entry.startsWith('.') && entry !== 'assets')
    .map((entry) => join(fullDir, entry))
    .sort((a, b) => basename(a).localeCompare(basename(b), 'zh-Hans-CN', { numeric: true }))
    .flatMap((entry) => {
      const stats = statSync(entry)

      if (stats.isDirectory()) {
        const items = directoryItems(relative(root, entry))
        if (items.length === 0) return []

        const readme = join(entry, 'README.md')
        return [
          {
            text: existsSync(readme) ? titleFromMarkdown(readme) : basename(entry),
            collapsed: true,
            items
          }
        ]
      }

      if (!entry.endsWith('.md')) return []

      return [
        {
          text: titleFromMarkdown(entry),
          link: linkFromFile(entry)
        }
      ]
    })
}

function sectionSidebar(dir: string): DefaultTheme.SidebarItem[] {
  return [
    { text: '方向首页', link: `/${dir}/README` },
    { text: '路线图', link: `/${dir}/roadmap` },
    {
      text: '学习笔记',
      collapsed: false,
      items: directoryItems(`${dir}/notes`)
    },
    {
      text: '每日记录',
      collapsed: true,
      items: directoryItems(`${dir}/daily`)
    }
  ]
}

const sidebar: DefaultTheme.Sidebar = Object.fromEntries(
  sections.map((section) => [`/${section.dir}/`, sectionSidebar(section.dir)])
)

export default defineConfig({
  title: '狸力',
  description: 'Software Product Engineer Roadmap',
  lang: 'zh-CN',
  base: '/Lili/',
  cleanUrls: true,
  ignoreDeadLinks: [/^(?:\.\/)?\.\.\/\.\.\/examples\//],
  lastUpdated: true,
  markdown: {
    lineNumbers: true
  },
  vite: {
    plugins: [
      {
        name: 'copy-roadmap-examples',
        closeBundle() {
          for (const section of sections) {
            const source = join(root, section.dir, 'examples')
            if (!existsSync(source)) continue

            cpSync(source, join(root, '.vitepress/dist', section.dir, 'examples'), {
              recursive: true
            })
          }
        }
      }
    ]
  },
  themeConfig: {
    logo: '/favicon.svg',
    nav: [
      { text: '首页', link: '/' },
      { text: '学习笔记', link: '/learning-notes' },
      ...sections.map((section) => ({
        text: section.text,
        link: `/${section.dir}/README`
      })),
      { text: 'GitHub', link: repo }
    ],
    sidebar: {
      '/': [
        { text: '总览', link: '/' },
        { text: '学习笔记知识库', link: '/learning-notes' },
        { text: '更新记录', link: '/CHANGELOG' },
        {
          text: '五个方向',
          items: sections.map((section) => ({
            text: section.text,
            link: `/${section.dir}/README`
          }))
        }
      ],
      ...sidebar
    },
    socialLinks: [{ icon: 'github', link: repo }],
    search: {
      provider: 'local'
    },
    outline: {
      level: [2, 3],
      label: '本文目录'
    },
    docFooter: {
      prev: '上一篇',
      next: '下一篇'
    },
    lastUpdated: {
      text: '最后更新',
      formatOptions: {
        dateStyle: 'medium',
        timeStyle: 'short'
      }
    },
    editLink: {
      pattern: `${repo}/edit/main/:path`,
      text: '在 GitHub 上编辑此页'
    },
    footer: {
      message: '从前端走向软件产品工程的长期路线图。',
      copyright: 'Copyright © justCDQ'
    }
  }
})
