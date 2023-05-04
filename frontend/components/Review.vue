<template>
  <div class="">
    <div class="flex flex-col lg:flex-row p-3 mx-auto place-content-center">
      <div class="flex-none text-md lg:text-5xl font-extrabold text-center lg:text-right py-3 xl:max-w-2xl md:max-w-md lg:max-w-xl">
        <span>(#{{ model.pr.number }})</span>
        <GradientText>{{ model.pr.title }}</GradientText>
      </div>
      <div class="flex-none xl:max-w-prose md:max-w-md lg:max-w-lg prose">
        <div class="md:m-3 p-3 bg-gray-50 rounded h-full border">
          <div
            v-if="model.pr.body"
            v-html="prBody"
          />
          <div
            v-if="model.review.body"
            v-html="reviewBody"
          />
        </div>
      </div>
    </div>
    <div class="gap-3">
      <Comment
        v-for="(comment, idx) in model.comments"
        :comment="comment"
        :idx="idx + 1"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import showdown from 'showdown';

const converter = new showdown.Converter();
converter.setFlavor('github');

const props = defineProps({
  model: { type: Object, required: true }
});

const prBody = computed(() => {
  return converter.makeHtml(props.model.pr.body);
});

const reviewBody = computed(() => {
  return converter.makeHtml(props.model.review.body);
});
</script>
