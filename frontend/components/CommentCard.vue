<template>
  <ComponentCard :badge="idx">
    <template #title>
      <code>{{ comment.path }}</code>
    </template>
    <template #left>
      <MarkdownHTML>{{ commentBody }}</MarkdownHTML>
    </template>
    <template #right>
      <DiffBlock :comment="comment" />
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
