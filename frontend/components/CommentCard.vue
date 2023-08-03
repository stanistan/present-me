<template>
  <ComponentCard :badge="comment.number">
    <template #title>
      <MaybeLink :href="comment.title.href" target="_blank">
        <code>{{ comment.title.text }}</code>
      </MaybeLink>
    </template>
    <template #body>
      <div class="flex flex-col md:flex-row max-h-[95vh] bg-gray-50">
        <div class="p-3 flex-none md:w-2/5 text-md markdown">
          <MarkdownHTML>{{ comment.description }}</MarkdownHTML>
        </div>
        <div
          class="flex-grow overflow-scroll text-sm md:border-l border-t md:border-t-0"
        >
          <DiffBlock :code="comment.code" />
        </div>
      </div>
    </template>
  </ComponentCard>
</template>

<script setup lang="ts">
import { Comment } from "../src/Review";
defineProps<{
  comment: Comment;
}>();
</script>
