<template>
  <div>
    <ComponentCard :badge="idx" title-class="">
      <template #title>
        <code>{{ comment.title.text }}</code>
      </template>
      <template #body>
        <div class="overflow-scroll text-lg max-h-[70vh]">
          <DiffBlock :code="comment.code" />
        </div>
      </template>
    </ComponentCard>
    <div class="max-w-[80%] mx-auto markdown">
      <MarkdownHTML>{{ commentBody }}</MarkdownHTML>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  comment: { type: Object, required: true },
  idx: { type: Number, required: true },
});

const commentBody = computed(() => {
  return props.comment.description.replace(/^\s*\d+\.\s*/, "");
});
</script>
