import fields from '@/server/data/fields.json'

// const fields = {}
export default defineEventHandler((event) => {
  const gvk = decodeURIComponent(`${event.context.params.gvk}`)
  return fields[gvk as keyof typeof fields]
})
