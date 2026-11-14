import { useCallback, useEffect, useRef, useState } from 'react'
import { useLang } from './i18n'
import { MAP_CENTER, MAP_ZOOM, ROUTE_URL } from './config/contacts'
import { IconChevronRight, IconAlertCircle } from './icons'

interface YmapsApi {
  ready: (callback: () => void) => void
  Map: new (
    container: HTMLElement,
    state: { center: [number, number]; zoom: number; controls?: string[] },
  ) => {
    geoObjects: { add: (object: unknown) => void }
    destroy: () => void
  }
  Placemark: new (
    coords: [number, number],
    properties: Record<string, string>,
    options?: Record<string, unknown>,
  ) => unknown
}

declare global {
  interface Window {
    ymaps?: YmapsApi
  }
}

const MAPS_API_BASE = 'https://api-maps.yandex.ru/2.1/'

let mapsLoader: Promise<YmapsApi> | null = null

function loadYmaps(apiKey: string): Promise<YmapsApi> {
  if (window.ymaps) return Promise.resolve(window.ymaps)
  if (mapsLoader) return mapsLoader

  mapsLoader = new Promise<YmapsApi>((resolve, reject) => {
    if (!apiKey) {
      reject(new Error('YANDEX_MAPS key is not configured'))
      return
    }
    const script = document.createElement('script')
    script.src = `${MAPS_API_BASE}?apikey=${encodeURIComponent(apiKey)}&lang=ru_RU`
    script.async = true
    script.addEventListener('load', () => {
      const api = window.ymaps
      if (!api) {
        reject(new Error('ymaps is missing'))
        return
      }
      api.ready(() => resolve(api))
    })
    script.addEventListener('error', () => {
      mapsLoader = null
      reject(new Error('Failed to load Yandex Maps script'))
    })
    document.head.appendChild(script)
  })
  return mapsLoader
}

const MAPS_API_KEY = (import.meta.env.VITE_YANDEX_MAPS ?? '') as string

export function YandexMap() {
  const { t } = useLang()
  const containerRef = useRef<HTMLDivElement>(null)
  const mapRef = useRef<InstanceType<YmapsApi['Map']> | null>(null)
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading')
  const [attempt, setAttempt] = useState(0)

  useEffect(() => {
    let cancelled = false
    setStatus('loading')

    loadYmaps(MAPS_API_KEY)
      .then((ymaps) => {
        if (cancelled || !containerRef.current) return
        const map = new ymaps.Map(containerRef.current, {
          center: MAP_CENTER,
          zoom: MAP_ZOOM,
          controls: ['zoomControl'],
        })
        const placemark = new ymaps.Placemark(
          MAP_CENTER,
          {
            balloonContentHeader: t('org.full'),
            balloonContentBody: t('contacts.address'),
          },
          { preset: 'islands#blueDotIcon', iconColor: '#1d4ed8' },
        )
        map.geoObjects.add(placemark)
        mapRef.current = map
        setStatus('ready')
      })
      .catch(() => {
        if (!cancelled) setStatus('error')
      })

    return () => {
      cancelled = true
      mapRef.current?.destroy()
      mapRef.current = null
    }
  }, [attempt, t])

  const retry = useCallback(() => setAttempt((n) => n + 1), [])

  return (
    <div className="map-block">
      <div className="map-frame">
        <div className="ymap" ref={containerRef} />

        {status === 'loading' ? (
          <div className="map-overlay">
            <span className="spinner" aria-hidden="true" />
            {t('contacts.mapLoading')}
          </div>
        ) : null}

        {status === 'error' ? (
          <div className="map-overlay map-overlay-error" role="alert">
            <IconAlertCircle />
            <span>{t('contacts.mapError')}</span>
            <button type="button" className="btn btn-ghost" onClick={retry}>
              {t('contacts.mapRetry')}
            </button>
          </div>
        ) : null}
      </div>

      <a className="btn btn-primary map-route" href={ROUTE_URL} target="_blank" rel="noreferrer">
        {t('contacts.routeButton')}
        <IconChevronRight width={16} height={16} />
      </a>
    </div>
  )
}
