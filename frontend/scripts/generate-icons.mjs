import { execFileSync } from 'node:child_process'
import { writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const outPath = join(here, '..', 'src', 'icons.tsx')

const ICONS = [
  ['IconSend', 'lucide:send'],
  ['IconFileText', 'lucide:file-text'],
  ['IconPaperclip', 'lucide:paperclip'],
  ['IconLogOut', 'lucide:log-out'],
  ['IconChevronLeft', 'lucide:chevron-left'],
  ['IconChevronRight', 'lucide:chevron-right'],
  ['IconX', 'lucide:x'],
  ['IconDownload', 'lucide:download'],
  ['IconInbox', 'lucide:inbox'],
  ['IconAlertCircle', 'lucide:alert-circle'],
  ['IconArrowLeft', 'lucide:arrow-left'],
  ['IconTrash', 'lucide:trash-2'],
  ['IconLoader', 'lucide:loader-circle'],
  ['IconSun', 'lucide:sun'],
  ['IconMoon', 'lucide:moon'],
]

function fetchInnerSvg(iconId) {
  const raw = execFileSync('npx', ['better-icons', 'get', iconId, '--json'], {
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'inherit'],
  })
  const { svg } = JSON.parse(raw)
  const start = svg.indexOf('>') + 1
  const end = svg.lastIndexOf('</svg>')
  let inner = svg.slice(start, end)
  inner = inner.replaceAll('stroke-width="2"', 'stroke-width="1.5"')
  inner = inner.replaceAll(' fill="currentColor"', '')
  return inner
}

const components = []
const bodies = []

for (const [name, iconId] of ICONS) {
  const body = fetchInnerSvg(iconId)
  bodies.push(`  ${name}: '${body.replaceAll("'", "\\'")}',`)
  components.push(
    `export function ${name}(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.${name} }} />
}
`,
  )
}

const file = `/* Сгенерировано scripts/generate-icons.mjs — не редактировать вручную.
   Источник: better-icons (Iconify), набор lucide. */

import type { SVGProps } from 'react'

export type IconProps = Omit<SVGProps<SVGSVGElement>, 'children'>

const BODIES = {
${bodies.join('\n')}
} as const

function Svg(props: IconProps) {
  return (
    <svg
      width={18}
      height={18}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      {...props}
    />
  )
}

${components.join('\n')}
`

writeFileSync(outPath, file)
console.log(`OK: ${ICONS.length} иконок -> ${outPath}`)
