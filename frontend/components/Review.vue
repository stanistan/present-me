<template>
  <div class="">
    <div>
      <div>This is a review {{ model.params.owner }}</div>
      <div>{{ model.pr.body }}</div>
      <div>{{ model.review.body }}</div>
    </div>
    <div class="gap-3">
      <Comment v-for="(comment, idx) in model.comments" :comment="comment" :idx="idx + 1" />
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  model: { type: Object, required: true }
})
</script>
