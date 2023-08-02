<template>
  <div v-html="rendered" />
</template>

<script setup lang="ts">
import showdown from "showdown";

const converter = new showdown.Converter();
converter.setFlavor("github");

const slots = useSlots()!!;
const rendered = computed(() => {
  //
  // holy cow typescript hates this.
  const d = slots.default!!();
  if (d.length > 0) {
    return converter.makeHtml(d[0].children as string);
  }

  // TODO: error?
  return "N/A";
});
</script>
