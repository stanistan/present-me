<template>
  <div class="markdown" v-if="body.length">
    <MarkdownHTML>
    {{ body }}
    </MarkdownHTML>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  model: { type: Object, required: true }
});

const body = computed(() => {
 return [ props.model.pr.body, props.model.review.body ]
    .filter(x => x && x.length > 0)
    .join("\n\n---\n\n");
});
</script>
