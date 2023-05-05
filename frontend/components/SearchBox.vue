<template>
  <form class="my-4 mx-auto text-lg max-w-3xl " @submit.prevent="$emit('submit')">
    <ErrorMessage :message="errorMessage" />
    <div class="mx-2 flex flex-row rounded bg-white shadow-md p-2 gap-2 border border-violet-100">
      <input 
        :disabled="disabled" :value="modelValue" name="search"
        type="text" 
        placeholder="$org/$repo/pull/$pull#pullrequestreview-$review" 
        class="flex-grow px-4 font-mono focus:ring-none rounded overflow-hidden inline-block"
        @input="$emit('update:modelValue', $event.target.value)"
      >
      <button
        type="submit"
        :disabled="disabled"
        class="rounded p-4 px-6 text-lg font-bold bg-gradient-to-b from-purple-700 to-purple-800
                      hover:from-purple-600 hover:to-purple-700
                      active:from-purple-600 active:to-purple-300
                      disabled:from-purple-200 disabled:to-purple-300 disabled:border-gray-100
                      border border-gray-600 hover:border-gray-400
                      text-white shadow-md"
      >
        <span v-if="disabled" class="animate-pulse">
          ...
        </span>
        <span v-else>
          go
        </span>
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
defineProps({
  modelValue: { type: String, default: "" }, 
  errorMessage: { type: String, default: "" }, 
  disabled: { type: Boolean, default: false }
});
defineEmits(['update:modelValue', 'submit']);
</script>
