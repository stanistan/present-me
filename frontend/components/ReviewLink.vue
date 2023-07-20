<template>
  <NuxtLink class="hover:no-underline" :class="active" :href="href">
    {{ to }}
  </NuxtLink>
</template>

<script setup type="ts">
const props = defineProps({
  params: { type: Object, required: true },
  to: { type: String, required: true },
  current: { type: String, default: "" }
});

const href = computed(() => {
  const { params, to } = props;
  const permalink = `/${params.owner}/${params.repo}/pull/${params.number}/review-${params.review}`;
  return to === 'cards' ? permalink : `${permalink}/slides`;
});

const active = computed(() => {
  const { current, to } = props;
  return current == to ? "font-bold no-underline" : "underline";
});
</script>
