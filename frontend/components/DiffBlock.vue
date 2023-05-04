<template>
  <pre><code ref="code" class="language-diff">{{ diff }}</code></pre>
</template>

<script setup lang="ts">
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css';

const code = ref('code');
onMounted(() => {
  hljs.highlightElement(code.value);
});

const props = defineProps({ 
  content: { 
    type: String,
    default: "",
  }
});
const diff = computed(() => {
  return props.content.split("\n").slice(1).join("\n");
});

</script>
