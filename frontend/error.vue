<template>
  <NuxtLayout>
    <div class="mx-auto max-w-3xl py-10">
      <div class="bg-orange-100 px-2 pb-2">
        <div class="text-xs text-center underline py-4">
          {{ error }}
        </div>
        <div class="bg-white p-4 text-center border border-orange-200 rounded">
          <button class="underline hover:no-underline" @click="handleError">
            go home
          </button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
defineProps({
  error: { type: Object, required: true }
});

const handleError = () => clearError({ redirect: '/' })
</script>
