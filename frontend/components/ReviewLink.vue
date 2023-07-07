<template>
  <NuxtLink class="underline hover:no-underline" :href="href">
    {{ to }}
  </NuxtLink>
</template>

<script setup type="ts">
const props = defineProps({
  params: { type: Object, required: true },
  to: { type: String, required: true }
});

const href = computed(() => {
  const { params, to } = props;
  const permalink = `/${params.owner}/${params.repo}/pull/${params.number}/review-${params.review}`;
  return to === 'cards' ? permalink : `${permalink}/slides`;
});
</script>
