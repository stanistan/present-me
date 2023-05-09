<template>
  <ComponentCard :badge="idx">
    <template #title>
      <code>{{ comment.path }}</code>
    </template>
    <template #body>
      <div class="flex flex-col md:flex-row max-h-[95vh] bg-gray-50">
        <div class="p-3 flex-none md:w-2/5 text-md markdown">
          <MarkdownHTML>{{ commentBody }}</MarkdownHTML>
        </div>
        <div class="flex-grow overflow-scroll text-sm md:border-l border-t md:border-t-0">
          <DiffBlock :comment="comment" />
        </div>
      </div>
    </template>
  </ComponentCard>
</template>

<script setup lang="ts">
const props = defineProps({
  comment: { type: Object, required: true },
  idx: { type: Number, required: true }
});

const commentBody = computed(() => {
  return props.comment.body.replace(/^\s*\d+\.\s*/, '');
});
</script>
