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
    <div class="flex flex-col md:flex-row max-h-[95vh] bg-gray-50">
      <div class="p-3 flex-none md:w-2/5 text-md markdown">
        <slot name="left" />
      </div>
      <div class="flex-grow overflow-scroll text-sm md:border-l border-t md:border-t-0">
        <slot name="right" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  badge: { type: Number, required: true },
});
</script>
