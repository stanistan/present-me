<template>
  <pre class="bg-gray-100"><code ref="code" :class="language">{{ comment.diff_hunk }}</code></pre>
</template>

<script setup lang="ts">
// We do manual Prism-ing
import Prism from 'prismjs';
Prism.manual = true;

// these are the supported languages in general
import 'prismjs/components/prism-bash';
import 'prismjs/components/prism-css';
import 'prismjs/components/prism-diff';
import 'prismjs/components/prism-docker';
import 'prismjs/components/prism-go';
import 'prismjs/components/prism-json';
import 'prismjs/components/prism-markdown';
import 'prismjs/components/prism-markup-templating';
import 'prismjs/components/prism-php';
import 'prismjs/components/prism-rust';
import 'prismjs/components/prism-sass';
import 'prismjs/components/prism-scss';
import 'prismjs/components/prism-toml';
import 'prismjs/components/prism-yaml';

// this gets us the fancier diff-highlighting
import 'prismjs/plugins/diff-highlight/prism-diff-highlight';
import 'prismjs/plugins/diff-highlight/prism-diff-highlight.css';

const code = ref('code');
onMounted(() => {
  Prism.highlightElement(code.value);
});

const props = defineProps({ 
  comment: { 
    type: Object,
    required: true,
  },
});

// we grab the file extension and map it to the diff-language
const languageMap = { rs: 'rust', vue: 'html', Dockerfile: 'docker' };
const language = computed(() => {
  const pieces = props.comment.path.split('.');
  const lang = pieces[pieces.length - 1];
  return `diff-highlight language-diff-${languageMap[lang] || lang}`;
});
</script>

<style type="text/css">
/**
 * prism.js default theme for JavaScript, CSS and HTML
 * Based on dabblet (http://dabblet.com)
 * @author Lea Verou
 */

code[class*="language-"],
pre[class*="language-"] {
	color: black;
	background: none;
	text-shadow: 0 1px #fefefe;
	font-size: 1em;
	text-align: left;
	white-space: pre;
	word-spacing: normal;
	word-break: normal;
	word-wrap: normal;
	line-height: 1.5;

	-moz-tab-size: 4;
	-o-tab-size: 4;
	tab-size: 4;

	-webkit-hyphens: none;
	-moz-hyphens: none;
	-ms-hyphens: none;
	hyphens: none;
}

pre[class*="language-"]::-moz-selection, pre[class*="language-"] ::-moz-selection,
code[class*="language-"]::-moz-selection, code[class*="language-"] ::-moz-selection {
	text-shadow: none;
	background: #b3d4fc;
}

pre[class*="language-"]::selection, pre[class*="language-"] ::selection,
code[class*="language-"]::selection, code[class*="language-"] ::selection {
	text-shadow: none;
	background: #b3d4fc;
}

@media print {
	code[class*="language-"],
	pre[class*="language-"] {
		text-shadow: none;
	}
}

/* Code blocks */
pre[class*="language-"] {
/*	padding: 1em;
	margin: .5em 0;
	overflow: auto; */
}

:not(pre) > code[class*="language-"],
pre[class*="language-"] {
	/*background: #f5f2f0; */
}

/* Inline code */
:not(pre) > code[class*="language-"] {
/*	padding: .1em;
	border-radius: .3em; */
	white-space: normal;
}

.token.comment,
.token.prolog,
.token.doctype,
.token.cdata {
	color: slategray;
}

.token.punctuation {
	color: #999;
}

.token.namespace {
	opacity: .7;
}

.token.property,
.token.tag,
.token.boolean,
.token.number,
.token.constant,
.token.symbol,
.token.deleted {
	color: #905;
}

.token.selector,
.token.attr-name,
.token.string,
.token.char,
.token.builtin,
.token.inserted {
	color: #690;
}

.token.operator,
.token.entity,
.token.url,
.language-css .token.string,
.style .token.string {
	color: #9a6e3a;
	/* This background color was intended by the author of this theme. */
	background: hsla(0, 0%, 100%, .5);
}

.token.atrule,
.token.attr-value,
.token.keyword {
	color: #07a;
}

.token.function,
.token.class-name {
	color: #DD4A68;
}

.token.regex,
.token.important,
.token.variable {
	color: #e90;
}

.token.important,
.token.bold {
	font-weight: bold;
}
.token.italic {
	font-style: italic;
}

.token.entity {
	cursor: help;
}

.token.prefix.inserted, 
.token.prefix.deleted,
.token.prefix.unchanged {
  padding-right: 2em;
  font-size: 0.90em;
}

pre.diff-highlight > code .token.deleted:not(.prefix), pre > code.diff-highlight .token.deleted:not(.prefix) {
  background-color: rgba(255, 0, 0, 0.05);
}

pre.diff-highlight > code .token.inserted:not(.prefix), pre > code.diff-highlight .token.inserted:not(.prefix) {
  background-color: rgba(0, 255, 180, .05);
}
</style>


