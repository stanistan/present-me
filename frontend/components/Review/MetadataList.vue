<template>
  <div>
    <div
      v-for="m in model.links"
      :key="m.heading"
      class="grid grid-cols-2 gap-4 text-xs font-mono"
    >
      <div class="text-right p-1">
        <NuxtLink class="underline hover:no-underline" :href="m.href">
          {{ m.text }}
        </NuxtLink>
      </div>
      <div class="p-1">
        <strong>{{ m.label }}</strong>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  model: { type: Object, required: true },
});
</script>
