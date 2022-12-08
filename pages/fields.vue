<template>
  <div class="flex justify-left items-center m-4">
    <div>
      <div class="p-4">
        <a class="text-lg">Can I use</a>
        <input v-model="filter" placeholder="Kind, eg. Ingress" class="input input-bordered mx-5" />
        <a class="text-lg">?</a>
      </div>
      <div class="p-4">
        <ul>
          <li v-for="(gvk, index) in filteredGVKs" :key="index">
            <a :href="`/field/${encodeURIComponent(gvk)}`" target="_blank">
              {{ gvk }}
            </a>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  props: {
    name: String,
    msg: { type: String, required: true },
  },
  data() {
    return {
      filter: '',
      gvks: [] as string[],
    }
  },
  mounted() {
    const fetchKinds = async () => {
      const response = await fetch('/api/gvks')
      const data = (await response.json()) as string[]
      this.gvks = data
    }
    fetchKinds()
  },
  computed: {
    filteredGVKs() {
      if (!this.filter) {
        return this.gvks
      }
      return this.gvks.filter((kind) => kind.toLowerCase().includes(this.filter.toLowerCase()))
    },
  },
})
</script>
