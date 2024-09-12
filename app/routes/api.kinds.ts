import { json } from '@remix-run/node'
import type { LoaderFunction } from '@remix-run/node'

import gvks from '../../public/data/gvk_api_lifecycle.json'

export const loader: LoaderFunction = async () => {
  const kindSet = new Set<string>()
  for (let gv in gvks) {
    for (let k in gvks[gv as keyof typeof gvks]) {
      kindSet.add(k)
    }
  }
  return json(Array.from(kindSet).sort())
}
