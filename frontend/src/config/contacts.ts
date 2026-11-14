const env = import.meta.env

export const CONTACT_EMAIL = (env.VITE_CONTACT_EMAIL ?? '') as string

export const MAP_CENTER: [number, number] = [53.114786, 26.043941]
export const MAP_ZOOM = 16

export const ROUTE_URL =
  `https://yandex.by/maps/?rtext=~${MAP_CENTER[0]},${MAP_CENTER[1]}&rtt=auto`

export const DUTY_PHONE_TEL = '+375163484822'
export const RECEPTION_PHONE_TEL = '+375163484728'
