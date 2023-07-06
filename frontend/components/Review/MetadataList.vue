<template>
  <div>
    <div
      v-for="m in metadata" :key="m.heading"
      class="grid grid-cols-2 gap-4 text-xs font-mono"
    >
      <div class="text-right p-1">
        <NuxtLink class="underline hover:no-underline" :href="m.href">
          {{ m.text }}
        </NuxtLink>
      </div>
      <div class="p-1">
        <strong>{{ m.heading }}</strong> 
      </div>
    </div>
  </div>
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
      heading: "Post",
      text: permalink,
      href: permalink
    },
    { 
      heading: "Slides",
      text: `${permalink}/slides`,
      href: `${permalink}/slides`,
    }
  ];
});
</script>
