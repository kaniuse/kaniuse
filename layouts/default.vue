<template>
  <div>
    <div class="navbar bg-neutral text-neutral-content">
      <div class="flex-1">
        <a class="btn btn-ghost normal-case text-xl" href="/">Kaniuse</a>
        <div class="tabs tabs-boxed">
          <div class="tab" :class="{ 'tab-active': $route.path.startsWith('/kind') }">
            <NuxtLink href="/kinds">Kinds</NuxtLink>
          </div>
          <div class="tab" :class="{ 'tab-active': $route.path.startsWith('/field') }">
            <NuxtLink href="/fields">Fields</NuxtLink>
          </div>
        </div>
      </div>

      <div class="flex-none">
        <a class="btn btn-ghost normal-case text-xl" target="_blank" href="https://github.com/kaniuse/kaniuse"
          >View On GitHub</a
        >
      </div>
    </div>
    <div class="container mx-auto">
      <slot />
    </div>
  </div>
</template>

<script lang="ts">
export default defineComponent({})
</script>
