<template>
  <div v-html="rendered" />
</template>

<script setup lang="ts">
import showdown from "showdown";

const converter = new showdown.Converter();
converter.setFlavor("github");

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

  return converter.makeHtml(data[0].children as string);
});
</script>
