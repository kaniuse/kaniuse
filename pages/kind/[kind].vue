<template>
  <h1 class="text-3xl p-4">
    <a class="text-gray-500">Availability for </a><a>{{ $route.params.kind }}</a>
  </h1>
  <div class="flex">
    <div class="flex-none text-right p-4">
      <div v-for="item in gvks" class="p-2 h-10">
        <a>{{ item }}</a>
      </div>
      <div class="p-2 h-10">
        <a class="link" target="_blank" href="https://endoflife.date/kubernetes"> Kubernetes Version </a>
      </div>
    </div>
    <div class="flex-1 p-4">
      <div
        v-for="span in spans"
        class="grid h-10"
        :style="{
          gridTemplateColumns: 'repeat(' + versions.length + ', minmax(0, 1fr))',
        }"
      >
        <div
          v-for="(item, index) in span.spans"
          class="flex border-neutral"
          :class="{
            'border-l': index !== span.spans.length - 1,
            'border-x': index === span.spans.length - 1,
          }"
          :style="{
            gridColumn: versions.indexOf(item.start) + 1 + ' / ' + (versions.indexOf(item.end) + 2),
          }"
        >
          <div
            class="w-full h-4 my-auto"
            :class="{
              'bg-green-500': item.lifecycle === 'stable',
              'bg-yellow-500': item.lifecycle === 'deprecated',
            }"
          ></div>
        </div>
      </div>

      <div
        class="h-10 grid text-center divide-x divide-neutral"
        :style="{
          gridTemplateColumns: 'repeat(' + versions.length + ', minmax(0, 1fr))',
        }"
      >
        <div class="p-2" v-for="item in versions">
          <a>{{ item }}</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import semver from 'semver'

type GVK = {
  group: string
  version: string
  kind: string
}

type GVKAndOccurrence = {
  GVK: GVK
  Lifecycles: { kubernetesMinorRelease: string; APILifecycle: string }[]
}

type GVKSpan = {
  GVK: GVK
  spans: { start: string; end: string; lifecycle: string }[]
}

export default defineNuxtComponent({
  data() {
    return {
      occurrences: [] as GVKAndOccurrence[],
    }
  },
  mounted() {
    const fetchKind = async () => {
      const response = await fetch(`/api/kind/${this.$route.params.kind}`)
      const gvks = (await response.json()) as GVKAndOccurrence[]
      this.occurrences = gvks
    }
    fetchKind()
  },
  computed: {
    gvks(): string[] {
      const result = new Array<string>()
      for (const gvk of this.occurrences) {
        result.push(`${gvk.GVK.group}/${gvk.GVK.version} ${gvk.GVK.kind}`)
      }
      return result
    },
    versions(): string[] {
      const versions = new Set<string>()
      for (const gvk of this.occurrences) {
        for (const lifecycle of gvk.Lifecycles) {
          versions.add(lifecycle.kubernetesMinorRelease)
        }
      }
      return Array.from(versions).sort((a, b) => {
        const semverA = semver.coerce(a)
        const semverB = semver.coerce(b)
        return semverA!.compare(semverB!)
      })
    },
    spans(): GVKSpan[] {
      const result = new Array<GVKSpan>()
      for (const gvk of this.occurrences) {
        const spans = new Array<{ start: string; end: string; lifecycle: string }>()
        // aggregate spans with same lifecycle
        var start = ''
        var lifecycle = ''
        for (const [index, item] of gvk.Lifecycles.entries()) {
          if (index == 0) {
            start = item.kubernetesMinorRelease
            lifecycle = item.APILifecycle
          }
          if (lifecycle !== item.APILifecycle) {
            spans.push({ start: start, end: gvk.Lifecycles[index - 1].kubernetesMinorRelease, lifecycle: lifecycle })
            start = item.kubernetesMinorRelease
            lifecycle = item.APILifecycle
          }
          if (index == gvk.Lifecycles.length - 1) {
            spans.push({ start: start, end: item.kubernetesMinorRelease, lifecycle: lifecycle })
          }
        }
        result.push({ GVK: gvk.GVK, spans: spans })
      }
      return result
    },
  },
})
</script>
