<template>
  <pre><code ref="code" class="diff-highlight" :class="language">{{ diff }}</code></pre>
</template>

<script setup lang="ts">
import Prism from 'prismjs';
Prism.manual = true;

import 'prismjs/components/prism-bash';
import 'prismjs/components/prism-toml';
import 'prismjs/components/prism-diff';
import 'prismjs/components/prism-json';
import 'prismjs/components/prism-markdown';
import 'prismjs/components/prism-markup-templating';
import 'prismjs/components/prism-php';
import 'prismjs/components/prism-scss';
import 'prismjs/components/prism-rust';
import 'prismjs/components/prism-go';

import 'prismjs/plugins/diff-highlight/prism-diff-highlight';
import 'prismjs/plugins/diff-highlight/prism-diff-highlight.css';

import 'prismjs/themes/prism.css';

const code = ref('code');
onMounted(() => {
  Prism.highlightElement(code.value);
});

const props = defineProps({ 
  content: { 
    type: String,
    default: ""
  },
  filename: {
    type: String,
    default: ""
  }
});

const diff = computed(() => {
  return props.content.split("\n").slice(1).join("\n");
});

const languageMap = {
  rs: 'rust',
};

const language = computed(() => {
  const pieces = props.filename.split('.');
  const lang = pieces[pieces.length - 1];
  return `language-diff-${languageMap[lang] || lang}`;
});

</script>
