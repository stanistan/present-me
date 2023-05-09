<template>
  <div class="m-4 border border-slate-300 rounded-xl overflow-hidden shadow">
    <div class="w-full">
      <div class="w-full text-sm p-3 bg-slate-100 border-b border-slate-300 gap-1 rounded-t-xl">
        <SmallBadge v-if="!!badge">
          {{ badge }}
        </SmallBadge>
        <slot name="title" />
      </div>
    </div>
    <slot class="bg-gray-50" name="body" />
  </div>
</template>

<script setup lang="ts">
defineProps({
  badge: { type: Number, default: 0 },
});
</script>
