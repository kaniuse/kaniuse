import gvks from "../../../../data/gvk_api_lifecycle.json";

export default defineEventHandler((event) => {
  const gv = `${event.context.params.group}/${event.context.params.version}`;
  return gvks[gv][event.context.params.kind];
});
