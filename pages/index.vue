<template>
  <div class="flex justify-center">
    <div class="container flex justify-left items-center m-4">
      <div>
        <a class="text-lg">Can I use</a>
        <input v-model="filter" placeholder="Type here" class="input input-bordered" />
        <a class="text-lg">?</a>

        <ul>
          <li v-for="item in filteredGvks">
            <a :href="`/api-lifecycle/${item.group}/${item.version}/${item.kind}`" target="_blank"
              >{{ item.group }}/{{ item.version }} - {{ item.kind }}</a
            >
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
      gvks: [] as {
        group: string
        version: string
        kind: string
      }[],
    }
  },
  mounted() {
    const fetchData = async () => {
      const response = await fetch('/api/gvks')
      const data = (await response.json()) as string[]
      data.map((item) => {
        const [group, versionkind] = item.split('/')
        const [version, kind] = versionkind.split('-')
        this.gvks.push({
          group: group.trim(),
          version: version.trim(),
          kind: kind.trim(),
        })
      })
    }
    fetchData()
  },
  computed: {
    filteredGvks() {
      return this.gvks.filter((item) => {
        return (
          item.group.toLowerCase().includes(this.filter.toLowerCase()) ||
          item.version.toLowerCase().includes(this.filter.toLowerCase()) ||
          item.kind.toLowerCase().includes(this.filter.toLowerCase())
        )
      })
    },
  },
})
</script>
