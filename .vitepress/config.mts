import { defineConfig, type DefaultTheme } from 'vitepress'
import { cpSync, existsSync, mkdirSync, readdirSync, readFileSync, statSync, writeFileSync } from 'node:fs'
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

const noteDirectoryRedirects: Record<string, string> = {
  '01-frontend/notes/web-basics': '01-frontend/notes/00-web-basics',
  '01-frontend/notes/html': '01-frontend/notes/01-html',
  '01-frontend/notes/css': '01-frontend/notes/02-css',
  '01-frontend/notes/javascript': '01-frontend/notes/03-javascript',
  '01-frontend/notes/typescript-frameworks': '01-frontend/notes/04-typescript-frameworks',
  '01-frontend/notes/browser-runtime': '01-frontend/notes/05-browser-runtime',
  '01-frontend/notes/application-architecture': '01-frontend/notes/06-application-architecture',
  '02-product/notes/foundations': '02-product/notes/00-foundations',
  '02-product/notes/product-teardowns': '02-product/notes/01-product-teardowns',
  '02-product/notes/problem-evidence': '02-product/notes/02-problem-evidence',
  '02-product/notes/requirements-prioritization': '02-product/notes/03-requirements-prioritization',
  '02-product/notes/solution-mvp': '02-product/notes/04-solution-mvp',
  '02-product/notes/prd-expression': '02-product/notes/05-prd-expression',
  '02-product/notes/metrics-experiments': '02-product/notes/06-metrics-experiments',
  '03-interaction-design/notes/foundations': '03-interaction-design/notes/00-foundations',
  '03-interaction-design/notes/principles': '03-interaction-design/notes/01-principles',
  '03-interaction-design/notes/deconstruction': '03-interaction-design/notes/02-deconstruction',
  '03-interaction-design/notes/validation': '03-interaction-design/notes/03-validation',
  '03-interaction-design/notes/information-architecture': '03-interaction-design/notes/04-information-architecture',
  '03-interaction-design/notes/flows-states': '03-interaction-design/notes/05-flows-states',
  '03-interaction-design/notes/interaction-patterns': '03-interaction-design/notes/06-interaction-patterns',
  '03-interaction-design/notes/06-interaction-patterns/navigation': '03-interaction-design/notes/06-interaction-patterns/00-navigation',
  '03-interaction-design/notes/06-interaction-patterns/input': '03-interaction-design/notes/06-interaction-patterns/01-input',
  '03-interaction-design/notes/06-interaction-patterns/data': '03-interaction-design/notes/06-interaction-patterns/02-data',
  '03-interaction-design/notes/06-interaction-patterns/feedback': '03-interaction-design/notes/06-interaction-patterns/03-feedback',
  '03-interaction-design/notes/06-interaction-patterns/operations': '03-interaction-design/notes/06-interaction-patterns/04-operations',
  '03-interaction-design/notes/06-interaction-patterns/collaboration': '03-interaction-design/notes/06-interaction-patterns/05-collaboration',
  '04-ai/notes/foundations': '04-ai/notes/00-foundations',
  '04-ai/notes/model-api': '04-ai/notes/01-model-api',
  '04-ai/notes/prompt': '04-ai/notes/02-prompt',
  '04-ai/notes/context-engineering': '04-ai/notes/03-context-engineering',
  '04-ai/notes/ai-ux': '04-ai/notes/04-ai-ux',
  '04-ai/notes/rag-parsing': '04-ai/notes/05-rag-parsing',
  '04-ai/notes/rag-chunking': '04-ai/notes/06-rag-chunking',
  '04-ai/notes/rag-retrieval': '04-ai/notes/07-rag-retrieval',
  '04-ai/notes/rag-evaluation': '04-ai/notes/08-rag-evaluation',
  '04-ai/notes/tool-design': '04-ai/notes/09-tool-design',
  '04-ai/notes/mcp': '04-ai/notes/10-mcp',
  '04-ai/notes/workflow': '04-ai/notes/11-workflow',
  '04-ai/notes/agent': '04-ai/notes/12-agent',
  '04-ai/notes/evaluation': '04-ai/notes/13-evaluation',
  '05-backend-data/notes/programming-basics': '05-backend-data/notes/00-programming-basics',
  '05-backend-data/notes/computer-systems': '05-backend-data/notes/01-computer-systems',
  '05-backend-data/notes/service-data-basics': '05-backend-data/notes/02-service-data-basics',
  '05-backend-data/notes/algorithms': '05-backend-data/notes/03-algorithms',
  '05-backend-data/notes/go': '05-backend-data/notes/04-go',
  '05-backend-data/notes/linux-network': '05-backend-data/notes/05-linux-network',
  '05-backend-data/notes/api-database': '05-backend-data/notes/06-api-database',
  '05-backend-data/notes/cache-messaging-storage': '05-backend-data/notes/07-cache-messaging-storage'
}

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

function writeRedirect(fromFile: string, toFile: string) {
  const relativeTarget = relative(dirname(fromFile), toFile).replaceAll('\\', '/')
  const href = relativeTarget.startsWith('.') ? relativeTarget : `./${relativeTarget}`

  mkdirSync(dirname(fromFile), { recursive: true })
  writeFileSync(
    fromFile,
    `<!doctype html><meta charset="utf-8"><meta http-equiv="refresh" content="0;url=${href}"><link rel="canonical" href="${href}"><script>location.replace(${JSON.stringify(href)})</script>`,
    'utf8'
  )
}

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
        name: 'copy-examples-and-note-redirects',
        closeBundle() {
          const dist = join(root, '.vitepress/dist')

          for (const section of sections) {
            const source = join(root, section.dir, 'examples')
            if (!existsSync(source)) continue

            cpSync(source, join(dist, section.dir, 'examples'), {
              recursive: true
            })
          }

          for (const [fromDir, toDir] of Object.entries(noteDirectoryRedirects)) {
            const targetDir = join(root, toDir)
            if (!existsSync(targetDir)) continue

            for (const entry of readdirSync(targetDir, { withFileTypes: true })) {
              if (!entry.isFile() || !entry.name.endsWith('.md')) continue

              const name = entry.name.replace(/\.md$/, '.html')
              writeRedirect(join(dist, fromDir, name), join(dist, toDir, name))
            }
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
