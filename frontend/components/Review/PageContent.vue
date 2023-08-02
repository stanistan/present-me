<template>
  <div class="gap-3">
    <div class="pt-4">
      <ReviewMetadataList :model="model" />
    </div>
    <ComponentCard>
      <template #title>
        <div class="text-xl font-extrabold">
          <GradientText>{{ model.title.text }}</GradientText>
        </div>
      </template>
      <template #body>
        <div class="px-4 py-4">
          <div v-if="model.body.length" class="markdown">
            <MarkdownHTML>
              {{ model.body }}
            </MarkdownHTML>
          </div>
        </div>
      </template>
    </ComponentCard>
    <CommentCard
      v-for="(comment, idx) in model.comments"
      :key="idx"
      :comment="comment"
    />
  </div>
</template>

<script setup lang="ts">
import { Review } from "../../src/Review";
defineProps<{
  model: Review;
}>();
</script>
