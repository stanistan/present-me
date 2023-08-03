<template>
  <NuxtLink
    v-if="href"
    :to="href"
    :target="target"
    class="underline hover:no-underline"
  >
    <slot />
  </NuxtLink>
  <slot v-else />
</template>

<script setup lang="ts">
defineProps<{
  target?: string;
  href?: string;
}>();
</script>
