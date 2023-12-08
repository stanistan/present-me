<template>
  <ReviewPageChrome v-slot="data" name="cards">
    <div class="mx-auto max-w-screen-2xl">
      <ReviewPageContent :model="data.data" />
    </div>
  </ReviewPageChrome>
</template>

<script setup lang="ts">
definePageMeta({
  name: "pull-tag"
});
</script>
