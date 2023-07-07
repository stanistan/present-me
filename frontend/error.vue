<template>
  <NuxtLayout>
    <div class="mx-auto max-w-3xl py-10">
      <div class="bg-orange-100 px-2 pb-2">
        <div class="text-xs text-center underline py-4">
          {{ error }}
        </div>
        <div class="bg-white p-4 text-center border border-orange-200 rounded">
          <button @click="handleError" class="underline hover:no-underline">go home</button>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
const props = defineProps({
  error: Object
})

const handleError = () => clearError({ redirect: '/' })
</script>
