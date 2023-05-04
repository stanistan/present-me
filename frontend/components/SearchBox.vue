<template>
  <form 
    @submit.prevent="$emit('submit')"
    class="my-4 mx-auto text-lg max-w-3xl ">
    <ErrorMessage :message="errorMessage" />
    <div class="mx-2 flex flex-row rounded bg-white shadow-md p-2 gap-2 border border-violet-100">
      <input 
         :disabled="disabled"
         :value="modelValue"
         @input="$emit('update:modelValue', $event.target.value)"
         name="search" 
         type="text" 
         placeholder="$org/$repo/pull/$pull#pullrequestreview-$review"
         class="flex-grow px-4 font-mono focus:ring-none rounded overflow-hidden inline-block" />
      <button
          type="submit"
          :disabled="disabled"
          class="rounded p-4 px-6 text-lg font-bold bg-gradient-to-b from-purple-700 to-purple-800
                      hover:from-purple-600 hover:to-purple-700
                      active:from-purple-600 active:to-purple-300
                      disabled:from-purple-200 disabled:to-purple-300 disabled:border-gray-100
                      border border-gray-600 hover:border-gray-400
                      text-white shadow-md">
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
const props = defineProps(['modelValue', 'errorMessage', 'disabled']);
const emits = defineEmits(['update:modelValue', 'submit']);


/*
async function executeSearch() {
  formDisabled.value = true;
  errorMessage.value = "";
  setTimeout(async function() {
    errorMessage.value = "failed with pretend.";
      formDisabled.value = false;

    const { data, error } = await useFetch('/api/search', {
      params: { search: props.searchQuery.value },
      server: false,
      initialCache: false,
      transform: v => JSON.parse(v)
    });

    if (error.value) {
      const errorData = JSON.parse(error.value.data)
      errorMessage.value = errorData.msg;
      formDisabled.value = false;
    } else {
      const params = data.value;
      await navigateTo(`${params.owner}/${params.repo}/pull/${params.number}/review-${params.review}`);
    }
  }, 2000);
}*/

</script>
