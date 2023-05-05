<template>
  <pre class="bg-gray-100"><code ref="code" :class="language">{{ diff }}</code></pre>
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

// drop the first line of the diff since it's a diff hunk header
const diff = computed(() => {
  const comment = props.comment;

  // let's get the lines first...
  const lines = comment.diff_hunk.split("\n");

  // the first line has metadata in it from the diff hunk...
  // and we can get the starting line that's expected to be there
  // from it...
  const startCountingAt = parseInt(lines.shift().match(/@@ -(\d+),/)[1], 10);

  // desiredRange 
  const 
    desiredStartLine = comment.start_line === undefined
      ? comment.line - 4
      : comment.start_line - (!startCountingAt ? 1 : 0),
    desiredEndLine = comment.line,
    outputDiff = [];

  let lineNumber = startCountingAt;
  lines.forEach(line => {
    if (lineNumber >= desiredStartLine && lineNumber <= desiredEndLine) {
      outputDiff.push(line);
    }

    if (comment.side == "LEFT" && !line.startsWith("+")) {
      lineNumber++;
    } else if (comment.side == "RIGHT" && !line.startsWith("-")) {
      lineNumber++;
    }
  });

  return outputDiff.join("\n");
});

// we grab the file extension and map it to the diff-language
const languageMap = { rs: 'rust' };
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


