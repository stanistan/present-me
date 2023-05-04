<template>
  <div class="rounded-lg font-bold ring-1 my-5 ring-red-300 bg-red-100 p-3 text-center" v-if="message">
    Error: <span class="underline">{{ message }}</span>
  </div>
</template>

<script setup lang="ts">
const props = defineProps([
  "message"
]);
</script>
