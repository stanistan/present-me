<template>
  <div v-html="rendered" />
</template>

<script setup lang="ts">
import { mdToHtml } from "../src/md";

const slots = useSlots();
const rendered = computed(() => {
  if (!slots) {
    return "";
  }
  
  if (!slots.default) {
    return "";
  }

  const data = slots.default();
  if (!data.length) {
    return "";
  }

  return mdToHtml(data[0].children as string);
});
</script>
