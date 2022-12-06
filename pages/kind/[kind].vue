<template>
  {{ $route.params.kind }}
  <div class="flex">
    <div class="flex-none text-right p-4">
      <div class="p-2 h-10"><a>networking.k8s.io/v1beta1 Ingress</a></div>
      <div class="p-2 h-10"><a>networking.k8s.io/v1 Ingress</a></div>
      <div class="p-2 h-10"><a>extensions/v1beta1 Ingress</a></div>
      <div class="p-2 h-10"><a>Kubernetes Versions</a></div>
    </div>
    <div class="flex-1 p-4">
      <div class="grid h-10" style="grid-template-columns: repeat(18, minmax(0, 1fr))">
        <div style="grid-column: span 6"></div>
        <div class="flex border-x border-r border-neutral" style="grid-column: span 8">
          <div class="bg-green-500 w-full h-4 my-auto" style=""></div>
        </div>
      </div>
      <div class="grid h-10" style="grid-template-columns: repeat(18, minmax(0, 1fr))">
        <div style="grid-column: span 11"></div>
        <div class="flex border-x border-neutral" style="grid-column: span 7">
          <div class="bg-green-500 w-full h-4 my-auto" style=""></div>
        </div>
      </div>
      <div class="grid h-10" style="grid-template-columns: repeat(18, minmax(0, 1fr))">
        <div class="flex border-l border-neutral" style="grid-column: span 6">
          <div class="bg-green-500 w-full h-4 my-auto" style=""></div>
        </div>
        <div class="flex border-x border-neutral" style="grid-column: span 8">
          <div class="bg-orange-500 w-full h-4 my-auto" style=""></div>
        </div>
      </div>
      <div
        class="h-10 grid text-center divide-x divide-neutral"
        style="grid-template-columns: repeat(18, minmax(0, 1fr))"
      >
        <div class="p-2"><a>1.8</a></div>
        <div class="p-2"><a>1.9</a></div>
        <div class="p-2"><a>1.10</a></div>
        <div class="p-2"><a>1.11</a></div>
        <div class="p-2"><a>1.12</a></div>
        <div class="p-2"><a>1.13</a></div>
        <div class="p-2"><a>1.14</a></div>
        <div class="p-2"><a>1.15</a></div>
        <div class="p-2"><a>1.16</a></div>
        <div class="p-2"><a>1.17</a></div>
        <div class="p-2"><a>1.18</a></div>
        <div class="p-2"><a>1.19</a></div>
        <div class="p-2"><a>1.20</a></div>
        <div class="p-2"><a>1.21</a></div>
        <div class="p-2"><a>1.22</a></div>
        <div class="p-2"><a>1.23</a></div>
        <div class="p-2"><a>1.24</a></div>
        <div class="p-2"><a>1.25</a></div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
type GVK = {
  Group: string
  Version: string
  Kind: string
}

type GVKAndOccurrence = {
  GVK: GVK
  Lifecycles: { kubernetesMinorRelease: string; APILifecycle: string }[]
}

export default defineNuxtComponent({
  data() {
    return {
      GVKs: [] as GVKAndOccurrence[],
    }
  },
  mounted() {
    const fetchKind = async () => {
      const response = await fetch(`/api/kind/${this.$route.params.kind}`)
      const gvks = (await response.json()) as GVKAndOccurrence[]
      this.GVKs = gvks
    }
    fetchKind()
  },
})
</script>
