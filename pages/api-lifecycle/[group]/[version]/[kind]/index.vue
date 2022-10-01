<template>
  {{ $route.params.group }}/{{ $route.params.version }}/{{ $route.params.kind }}
  <ul>
    <li v-for="item in apiLifecycles">
      <a>Kubernetes {{ item.kubernetesMinorRelease }}</a>
      <a>: {{ item.APILifecycle }}</a>
    </li>
  </ul>
</template>

<script lang="ts">
type APILifecycle = {
  kubernetesMinorRelease: string
  APILifecycle: string
}
export default defineNuxtComponent({
  data() {
    return {
      apiLifecycles: {} as APILifecycle[],
    }
  },
  mounted() {
    const fetchData = async () => {
      const response = await fetch(
        `/api/api-lifecycle/${this.$route.params.group}/${this.$route.params.version}/${this.$route.params.kind}`
      )
      const apiLifecycles = (await response.json()) as APILifecycle[]
      this.apiLifecycles = apiLifecycles
    }

    fetchData()
  },
})
</script>
