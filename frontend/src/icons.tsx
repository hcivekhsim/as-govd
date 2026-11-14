import type { SVGProps } from 'react'

export type IconProps = Omit<SVGProps<SVGSVGElement>, 'children'>

const BODIES = {
  IconSend: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11zm7.318-19.539l-10.94 10.939"/>',
  IconFileText: '<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path d="M6 22a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h8a2.4 2.4 0 0 1 1.704.706l3.588 3.588A2.4 2.4 0 0 1 20 8v12a2 2 0 0 1-2 2z"/><path d="M14 2v5a1 1 0 0 0 1 1h5M10 9H8m8 4H8m8 4H8"/></g>',
  IconPaperclip: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="m16 6l-8.414 8.586a2 2 0 0 0 2.829 2.829l8.414-8.586a4 4 0 1 0-5.657-5.657l-8.379 8.551a6 6 0 1 0 8.485 8.485l8.379-8.551"/>',
  IconLogOut: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="m16 17l5-5l-5-5m5 5H9m0 9H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>',
  IconChevronLeft: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="m15 18l-6-6l6-6"/>',
  IconChevronRight: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="m9 18l6-6l-6-6"/>',
  IconX: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M18 6L6 18M6 6l12 12"/>',
  IconDownload: '<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path d="M12 15V3m9 12v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><path d="m7 10l5 5l5-5"/></g>',
  IconInbox: '<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path d="M22 12h-6l-2 3h-4l-2-3H2"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11"/></g>',
  IconAlertCircle: '<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M12 8v4m0 4h.01"/></g>',
  IconArrowLeft: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="m12 19l-7-7l7-7m7 7H5"/>',
  IconTrash: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10 11v6m4-6v6m5-11v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>',
  IconLoader: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a9 9 0 1 1-6.219-8.56"/>',
  IconSun: '<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><circle cx="12" cy="12" r="4"/><path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/></g>',
  IconMoon: '<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20.985 12.486a9 9 0 1 1-9.473-9.472c.405-.022.617.46.402.803a6 6 0 0 0 8.268 8.268c.344-.215.825-.004.803.401"/>',
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

export function IconSend(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconSend }} />
}

export function IconFileText(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconFileText }} />
}

export function IconPaperclip(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconPaperclip }} />
}

export function IconLogOut(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconLogOut }} />
}

export function IconChevronLeft(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconChevronLeft }} />
}

export function IconChevronRight(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconChevronRight }} />
}

export function IconX(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconX }} />
}

export function IconDownload(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconDownload }} />
}

export function IconInbox(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconInbox }} />
}

export function IconAlertCircle(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconAlertCircle }} />
}

export function IconArrowLeft(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconArrowLeft }} />
}

export function IconTrash(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconTrash }} />
}

export function IconLoader(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconLoader }} />
}

export function IconSun(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconSun }} />
}

export function IconMoon(props: IconProps) {
  return <Svg {...props} dangerouslySetInnerHTML={{ __html: BODIES.IconMoon }} />
}

