<template>
  <NuxtLayout>
    <ErrorBlock>
      <template #title>
        {{ error }}
      </template>
      <button class="underline hover:no-underline" @click="handleError">
        go home
      </button>
    </ErrorBlock>
  </NuxtLayout>
</template>

<script setup lang="ts">
defineProps({
  error: { type: Object, required: true },
});

const handleError = () => clearError({ redirect: "/" });
</script>
