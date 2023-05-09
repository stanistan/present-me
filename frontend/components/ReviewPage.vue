<template>
  <div class="gap-3">
    <ComponentCard>
      <template #title>
        <div class="text-md text-xl font-extrabold">
          <span>(#{{ model.pr.number }})</span>&nbsp;
          <GradientText>{{ model.pr.title }}</GradientText>
        </div>
      </template>
      <template #left>
        <MarkdownHTML v-if="model.pr.body">
        {{ model.pr.body }}
        </MarkdownHTML>
      </template>
      <template #right>
        <div class="px-4 py-4">
          <ul class="list-disc px-4 mb-4 text-xs">
            <li class="prose">
              <strong>Author</strong> :: 
              <NuxtLink class="underline" :href="model.pr.user.html_url">{{ model.pr.user.login }}</NuxtLink>
            </li>
            <li>
              <strong>Source</strong> ::
              <NuxtLink class="underline" :href="model.pr.html_url">{{ model.pr.html_url }}</NuxtLink>
            </li>
            <li>
              <strong>Review by</strong> :: 
              <NuxtLink class="underline" :href="model.review.user.html_url">{{ model.review.user.login }}</NuxtLink>
            </li>
          </ul>
          <MarkdownHTML v-if="model.review.body">
          {{ model.review.body }}
          </MarkdownHTML>
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
