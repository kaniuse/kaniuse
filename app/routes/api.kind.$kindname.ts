import { json } from '@remix-run/node'
import { LoaderFunction } from '@remix-run/node'

import gvks from '../../public/data/gvk_api_lifecycle.json'

export const loader: LoaderFunction = async ({ params }) => {
  const kindname = params.kindname

  if (!kindname) {
    throw new Response('Kind name is required', { status: 400 })
  }

  const kindData = Object.entries(gvks).flatMap(([gv, kinds]) =>
    Object.entries(kinds)
      .filter(([k]) => k === kindname)
      .map(([k, lifecycles]) => ({
        GVK: {
          group: gv.split('/')[0],
          version: gv.split('/')[1],
          kind: k,
        },
        Lifecycles: lifecycles,
      }))
  )

  if (kindData.length === 0) {
    throw new Response('Not Found', { status: 404 })
  }

  return json(kindData)
}
