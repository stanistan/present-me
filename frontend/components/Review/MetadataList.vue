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
  return [
    { 
      heading: "Author", 
      text: pr.user.login, 
      href: pr.user.html_url,
    },
    { 
      heading: "Pull Request",
      text: `${params.owner}/${params.repo}/pull/${params.number}`,
      href: pr.html_url,
    },
    { 
      heading: "Review",
      text: `#review-${review.id}`,
      href: review.html_url,
    },
  ];
});
</script>
