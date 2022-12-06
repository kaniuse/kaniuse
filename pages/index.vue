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
          <li v-for="(kind, index) in filteredKinds" :key="index">
            <a :href="`/kind/${kind}`" target="_blank">
              {{ kind }}
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
      kinds: [] as string[],
    }
  },
  mounted() {
    const fetchKinds = async () => {
      const response = await fetch('/api/kinds')
      const data = (await response.json()) as string[]
      this.kinds = data
    }
    fetchKinds()
  },
  computed: {
    filteredKinds() {
      if (!this.filter) {
        return this.kinds
      }
      return this.kinds.filter((kind) => kind.toLowerCase().includes(this.filter.toLowerCase()))
    },
  },
})
</script>
