<template>
  <div class="">
    <div class="">
      <div class="text-5xl font-extrabold text-center py-3">
        <span class="bg-clip-text text-transparent bg-gradient-to-r underline from-pink-500 to-violet-500">
          {{ model.pr.title }}
        </span>
        (#{{ model.pr.number }})
      </div>
      <div>{{ model.pr.body }}</div>
      <div>{{ model.review.title }}</div>
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
