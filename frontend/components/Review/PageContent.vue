<template>
  <div class="gap-3">
    <div class="pt-4">
          <Review-MetadataList :model="model" />
    </div>
    <ComponentCard>
      <template #title>
        <div class="text-xl font-extrabold">
          <span>(#{{ model.pr.number }})</span>&nbsp;
          <GradientText>{{ model.pr.title }}</GradientText>
        </div>
      </template>
      <template #body>
        <div class="px-4 py-4">
          <BodyMarkdown :model="model" />
        </div>
      </template>
    </ComponentCard>
    <CommentCard
      v-for="(comment, idx) in model.comments"
      :key="idx"
      :comment="comment"
      :idx="idx + 1"
    />
  </div>
</template>

<script setup lang="ts">
defineProps({
  model: { type: Object, required: true }
});
</script>
