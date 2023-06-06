<template>
  <div class="gap-3">
    <ComponentCard>
      <template #title>
        <div class="text-xl font-extrabold">
          <span>(#{{ model.pr.number }})</span>&nbsp;
          <GradientText>{{ model.pr.title }}</GradientText>
        </div>
      </template>
      <template #body>
        <div class="px-4 py-4">
          <Review-MetadataList :model="model" />
          <div class="markdown">
            <MarkdownHTML v-if="body.length">
              {{ body }}
            </MarkdownHTML>
          </div>
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
const props = defineProps({
  model: { type: Object, required: true }
});

const body = computed(() => {
 return  [ props.model.pr.body, props.model.review.body ]
    .filter(x => x && x.length > 0)
    .join("\n\n---\n\n");
});
</script>
