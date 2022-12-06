import kinds from '@/server/data/kinds.json'

export default defineEventHandler((event) => {
  const kind = `${event.context.params.kind}`
  // if gvks contains key gv
  return kinds[kind as keyof typeof kinds]
})
