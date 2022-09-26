import gvks from '../data/gvk_api_lifecycle.json';

export default defineEventHandler((event) => {
  let result :Array<String> = [];

  for (let gv in gvks) {
    for(let k in gvks[gv]) {
      result.push(gv + " - " + k);
    }
  }
  return result;
})
