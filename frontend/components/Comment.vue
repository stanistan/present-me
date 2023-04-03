<template>
    <div class="m-4 border border-slate-300 rounded-xl overflow-hidden shadow flex-grow">
      <div class="w-full">
        <div class="w-full text-sm p-3 bg-slate-100 border-b border-slate-300 gap-1rounded-t-xl">
          <span class="text-center ring bg-indigo-300 rounded-3xl p-1 text-xs mr-2 ring-indigo-100 font-mono">{{ idx }}</span>
          <code>{{ comment.path }}</code>
        </div>
      </div>
      <div class="flex flex-row max-h-[95vh] bg-gray-50">
        <div class="p-3 flex-none w-2/5 text-md">
          <div v-html="commentBody"></div>
        </div>
        <div class="flex-grow overflow-scroll text-sm border-l">
          <pre class=""><code class="language-diff">{{ comment.diff_hunk }}</code></pre>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import showdown from 'showdown';

import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';

const converter = new showdown.Converter();
converter.setFlavor('github');

const props = defineProps({
  comment: { type: Object, required: true },
  idx: { type: String, required: true }
});

const commentBody = computed(() => {
  const replaced = props.comment.body.replace(/^\s*\d+\.\s*/, '');
  return converter.makeHtml(replaced);
});

onMounted(() => {
  hljs.highlightAll();
});
</script>
