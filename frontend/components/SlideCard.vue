<template>
  <div>
    <ComponentCard :badge="comment.number" title-class="">
      <template #title>
        <MaybeLink :href="comment.title.href" target="_blank">
          <code>{{ comment.title.text }}</code>
        </MaybeLink>
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
import { Comment } from "../src/Review";
defineProps<{
  comment: Comment;
}>();
</script>
