import gvks from '../data/gvk_api_lifecycle.json'

export default defineEventHandler((event) => {
  const kindSet = new Set<string>()
  for (let gv in gvks) {
    for (let k in gvks[gv as keyof typeof gvks]) {
      kindSet.add(k)
    }
  }
  return Array.from(kindSet).sort()
})
