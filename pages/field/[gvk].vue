<template>
  <h1 class="text-3xl p-4">
    <a class="text-gray-500">Fields Availability for </a><a>{{ $route.params.gvk }}</a>
  </h1>
  <div class="overflow-x-auto overflow-y-auto h-screen">
    <table class="table table-compact w-full">
      <thead>
        <tr style="position: sticky">
          <th class="border px-4 py-2">Field</th>
          <th v-for="version in versions" class="border px-4 py-2">{{ version }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="field in fields" :key="field">
          <td class="border px-4 py-2">{{ field }}</td>
          <td v-for="version in versions" class="border px-4 py-2">{{ aggregatedFields.get(field)!.get(version) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script lang="ts">
import semver from 'semver'

export default defineComponent({
  data() {
    return {
      fieldSummary: [] as {
        FieldPath: string
        LifeCycles: {
          KubernetesMinorRelease: string
          APILifecycle: string
        }[]
      }[],
    }
  },
  mounted() {
    const fetchFields = async () => {
      const response = await fetch(`/api/field/${encodeURIComponent(this.$route.params.gvk as string)}`)
      const data = (await response.json()) as {
        GVK: {
          group: string
          version: string
          kind: string
        }
        KindLifeCycles: {
          kubernetesMinorRelease: string
          APILifecycle: string
        }[]
        Summary: any
      }

      for (const key in data.Summary) {
        const item = data.Summary[key]! as {
          FieldPath: string
          LifeCycles: {
            KubernetesMinorRelease: string
            APILifecycle: string
          }[]
        }
        this.fieldSummary.push(item)
      }
    }

    fetchFields()
  },
  computed: {
    versions(): string[] {
      const versionSet = new Set<string>()
      Array.from(this.fieldSummary)
        .flatMap((field) => field.LifeCycles.map((lc) => lc.KubernetesMinorRelease))
        .forEach((version) => versionSet.add(version))
      console.log(versionSet)
      return Array.from(versionSet).sort((a, b) => {
        const semverA = semver.coerce(a)
        const semverB = semver.coerce(b)
        return semverA!.compare(semverB!)
      })
    },
    aggregatedFields(): Map<string, Map<string, string>> {
      // layer 1 key is field path, layer 2 key is version, value is lifecycle
      const result = new Map<string, Map<string, string>>()
      this.fieldSummary.forEach((field) => {
        const fieldPath = field.FieldPath
        const fieldMap = new Map<string, string>()
        field.LifeCycles.forEach((lc) => {
          fieldMap.set(lc.KubernetesMinorRelease, lc.APILifecycle)
        })
        result.set(fieldPath, fieldMap)
      })
      return result
    },
    fields(): string[] {
      return Array.from(this.aggregatedFields.keys())
    },
  },
})
</script>
