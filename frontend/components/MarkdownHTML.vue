<template>
  <div v-html="rendered" />
</template>

<script setup lang="ts">
import showdown from 'showdown';

const converter = new showdown.Converter();
converter.setFlavor('github');

const slots = useSlots();
const rendered = computed(() => {
  return converter.makeHtml(slots['default']()[0].children);
});
</script>
