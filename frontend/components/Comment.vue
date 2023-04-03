<template>
  <div class="bg-gray-500 p-3 m-3">
    <div class="bg-white p-3">
      <div class="text-lg font-bold">File Path: {{ comment.path }}</div>
    </div>
    <div class="flex flex-row max-h-[90vh]">
      <div class="bg-yellow-100 p-3 flex-none w-2/5">
        <div v-html="commentBody"></div>
      </div>
      <div class="bg-yellow-100 p-1 flex-grow overflow-scroll text-sm">
        <pre style="border:0"><code class="language-diff">{{ comment.diff_hunk }}</code></pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import showdown from 'showdown';

import Prism from 'prismjs';
import 'prismjs/components/prism-diff'
import 'prism-themes/themes/prism-material-oceanic.css'

const converter = new showdown.Converter();
converter.setFlavor('github');

const props = defineProps({
  comment: { type: Object, required: true }
});

const commentBody = computed(() => {
  const replaced = props.comment.body.replace(/^\s*\d+\.\s*/, '');
  return converter.makeHtml(replaced);
});

onMounted(() => {
  Prism.highlightAll();
});


</script>
