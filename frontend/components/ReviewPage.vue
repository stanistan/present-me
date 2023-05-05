<template>
  <div class="">
    <div class="flex flex-col lg:flex-row p-3 mx-auto place-content-center">
      <div class="flex-none text-md lg:text-5xl font-extrabold text-center lg:text-right py-3 xl:max-w-2xl md:max-w-md lg:max-w-xl">
        <span>(#{{ model.pr.number }})</span>
        <GradientText>{{ model.pr.title }}</GradientText>
      </div>
      <div class="flex-none xl:max-w-prose md:max-w-md lg:max-w-lg prose">
        <div class="md:m-3 p-3 bg-gray-50 rounded h-full border">
          <MarkdownHTML v-if="model.pr.body">
            {{ model.pr.body }}
          </MarkdownHTML>
          <MarkdownHTML v-if="model.review.body">
            {{ model.review.body }}
          </MarkdownHTML>
        </div>
      </div>
    </div>
    <div class="gap-3">
      <CommentCard
        v-for="(comment, idx) in model.comments.slice(0, 4)"
        :key="idx"
        :comment="comment"
        :idx="idx + 1"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  model: { type: Object, required: true }
});
</script>
