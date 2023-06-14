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

  const { model } = props;
  const { pr, review, params } = model;
  const permalink = `/${params.owner}/${params.repo}/pull/${params.number}/review-${params.review}`;
  return [
    { 
      heading: "Author", 
      text: pr.user.login, 
      href: pr.user.html_url,
    },
    { 
      heading: "Pull Request",
      text: pr.html_url, 
      href: pr.html_url,
    },
    { 
      heading: "Pull Request Review",
      text: review.html_url, 
      href: review.html_url,
    },
    { 
      heading: "Review Author",
      text: review.user.login, 
      href: review.user.html_url,
    },
    {
      heading: "Permalink",
      text: permalink,
      href: permalink
    },
    { 
      heading: "Slides Permalink",
      text: `${permalink}/slides`,
      href: `${permalink}/slides`,
    }
  ];
});
</script>
