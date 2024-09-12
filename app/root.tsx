import type { LinksFunction, MetaFunction } from '@remix-run/node'
import { Links, LiveReload, Meta, Outlet, Scripts, ScrollRestoration } from '@remix-run/react'
import { Link, useLocation } from '@remix-run/react'

import styles from './tailwind.css?url'

export const links: LinksFunction = () => [{ rel: 'stylesheet', href: styles }]

export const meta: MetaFunction = () => [
  { charset: 'utf-8' },
  { title: 'Kaniuse' },
  { name: 'viewport', content: 'width=device-width,initial-scale=1' },
]

export default function App() {
  const location = useLocation()

  return (
    <html lang="en">
      <head>
        <Meta />
        <Links />
      </head>
      <body>
        <div>
          <div className="navbar bg-neutral text-neutral-content">
            <div className="flex-1">
              <Link to="/" className="btn btn-ghost normal-case text-xl">
                Kaniuse
              </Link>
              <div className="tabs tabs-boxed">
                <div className={`tab ${location.pathname.startsWith('/kind') ? 'tab-active' : ''}`}>
                  <Link to="/kinds">Kinds</Link>
                </div>
                <div className={`tab ${location.pathname.startsWith('/field') ? 'tab-active' : ''}`}>
                  <Link to="/fields">Fields</Link>
                </div>
              </div>
            </div>

            <div className="flex-none">
              <a
                className="btn btn-ghost normal-case text-xl"
                target="_blank"
                href="https://github.com/kaniuse/kaniuse"
                rel="noopener noreferrer"
              >
                ⭐ Star On GitHub
              </a>
            </div>
          </div>
          <div className="container mx-auto h-full">
            <Outlet />
          </div>
        </div>
        <ScrollRestoration />
        <Scripts />
        <LiveReload />
      </body>
    </html>
  )
}
