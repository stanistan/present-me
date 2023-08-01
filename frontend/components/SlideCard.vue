<template>
  <div>
    <ComponentCard :badge="comment.number" title-class="">
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
      <MarkdownHTML>{{ comment.description }}</MarkdownHTML>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  comment: { type: Object, required: true },
});
</script>
