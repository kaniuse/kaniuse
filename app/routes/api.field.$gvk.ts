import { json } from '@remix-run/node'
import type { LoaderFunction } from '@remix-run/node'

import fields from '../../public/data/fields.json'

export const loader: LoaderFunction = async ({ params }) => {
  const gvk = decodeURIComponent(params.gvk as string)

  if (!(gvk in fields)) {
    throw new Response('GVK not found', { status: 404 })
  }

  return json(fields[gvk as keyof typeof fields])
}
