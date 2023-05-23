<template>
  <div class="text-center">
    redirecting...
  </div>
</template>
<script setup type="ts">
const { params } = useRoute();
await navigateTo(`/${params.org}/${params.repo}/pull/${params.pull}/review-${params.review}`);
</script>
