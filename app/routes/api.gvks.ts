import { json } from '@remix-run/node'

import gvks from '../../public/data/gvk_api_lifecycle.json'

export async function loader() {
  let result: string[] = []

  for (let gv in gvks) {
    for (let k in gvks[gv as keyof typeof gvks]) {
      result.push(`${gv} ${k}`)
    }
  }

  return json(result)
}
