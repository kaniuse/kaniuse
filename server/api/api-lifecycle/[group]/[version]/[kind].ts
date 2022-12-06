import gvks from '@/server/data/gvk_api_lifecycle.json'

export default defineEventHandler((event) => {
  const gv = `${event.context.params.group}/${event.context.params.version}`
  // if gvks contains key gv
  if (gvks[gv as keyof typeof gvks]) {
    const kinds = gvks[gv as keyof typeof gvks]
    return kinds[event.context.params.kind as keyof typeof kinds]
  }
  return {}
})
