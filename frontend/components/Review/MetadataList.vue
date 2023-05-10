<template>
  <ul class="list-disc px-4 mb-4 text-xs">
    <li v-for="m in metadata" :key="m.heading">
      <strong>{{ m.heading }}</strong> ::
      <NuxtLink class="underline" :href="m.href">
        {{ m.text }}
      </NuxtLink>
    </li>
  </ul>
</template>

<script setup lang="ts">
const props = defineProps({
  model: { type: Object, required: true }
});

const metadata = computed(() => {
  return [
    { 
      heading: "Author", 
      text: props.model.pr.user.login, 
      href: props.model.pr.user.html_url,
    },
    { 
      heading: "Pull Request",
      text: props.model.pr.html_url, 
      href: props.model.pr.hrml_url,
    },
    { 
      heading: "Pull Request Review",
      text: props.model.review.html_url, 
      href: props.model.review.hrml_url,
    },
    { 
      heading: "Review Author",
      text: props.model.review.user.login, 
      href: props.model.review.user.html_url,
    },
  ];
});
</script>
