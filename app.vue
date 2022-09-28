<script setup lang="ts">
const raw = await useFetch('/api/gvks')
const gvks = raw.data.value.map((item) => {
  const gandvk = (item as string).split('/')
  const vandk = gandvk[1].split('-')
  return {
    group: gandvk[0].trim(),
    version: vandk[0].trim(),
    kind: vandk[1].trim(),
  }
})
</script>
<template>
  <!-- Just test daisyUI. -->
  <div class="p-3">
    <div class="navbar bg-base-100 shadow-xl rounded-box">
      <a class="btn btn-ghost normal-case text-xl">daisyUI</a>
    </div>
  </div>
  <div>
    <li v-for="item in gvks">
      <a :href="`/api/api-lifecycle/${item.group}/${item.version}/${item.kind}`" target="_blank"
        >{{ item.group }}/{{ item.version }} - {{ item.kind }}</a
      >
    </li>
  </div>
</template>
