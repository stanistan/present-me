<template>
  <div class="">
    <div>
      <div>This is a review {{ model.params.owner }}</div>
      <div>{{ model.pr.body }}</div>
      <div>{{ model.review.body }}</div>
    </div>
    <div class="gap-3">
      <Comment v-for="comment in model.comments" :comment="comment" />
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  model: { type: Object, required: true }
})
</script>
