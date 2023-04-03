<template>
  <div class="bg-gray-500 p-3 m-3">
    <div class="bg-white p-3">
      <div class="text-lg font-bold">File Path: {{ comment.path }}</div>
    </div>
    <div class="flex flex-row max-h-[90vh]">
      <div class="bg-yellow-100 p-3 flex-1">
        <div v-html="commentBody"></div>
      </div>
      <div class="bg-yellow-200 p-3 flex-1 overflow-scroll">
        <code class="language-diff text-sm">
          {{ comment.diff_hunk }}
        </code>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import showdown from 'showdown';
const converter = new showdown.Converter();
converter.setFlavor('github');

const props = defineProps({
  comment: { type: Object, required: true }
});

const commentBody = computed(() => {
  const replaced = props.comment.body.replace(/^\s*\d+\.\s*/, '');
  return converter.makeHtml(replaced);
});

</script>
