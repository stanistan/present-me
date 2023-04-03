<template>
  <div class="">
    <div class="flex flex-row p-3 mx-auto place-content-center">
      <div class="flex-none text-5xl font-extrabold text-right py-3 max-w-2xl">
        <span class="bg-clip-text text-transparent bg-gradient-to-r underline from-pink-600 to-violet-900">
          {{ model.pr.title }}
        </span>
        (#{{ model.pr.number }})
      </div>
      <div class="flex-none max-w-prose prose">
        <div class="m-3 p-3 bg-gray-50 rounded h-full border">
          <div v-html="prBody"></div>
          <div v-html="reviewBody"></div>
        </div>
      </div>
    </div>
    <div class="gap-3">
      <Comment v-for="(comment, idx) in model.comments" :comment="comment" :idx="idx + 1" />
    </div>
  </div>
</template>

<script setup lang="ts">
import showdown from 'showdown';

const converter = new showdown.Converter();
converter.setFlavor('github');

const props = defineProps({
  model: { type: Object, required: true }
})

const prBody = computed(() => {
  return converter.makeHtml(props.model.pr.body);
})

const reviewBody = computed(() => {
  return converter.makeHtml(props.model.review.body);
})

</script>
