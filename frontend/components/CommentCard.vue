<template>
  <div class="m-4 border border-slate-300 rounded-xl overflow-hidden shadow">
    <div class="w-full">
      <div class="w-full text-sm p-3 bg-slate-100 border-b border-slate-300 gap-1 rounded-t-xl">
        <span class="text-center ring bg-gray-700 text-white rounded-3xl p-1 px-2 text-xs mr-2 ring-gray-100 font-mono">{{ idx }}</span>
        <code>{{ comment.path }}</code>
      </div>
    </div>
    <div class="flex flex-col md:flex-row max-h-[95vh] bg-gray-50">
      <div class="p-3 flex-none md:w-2/5 text-md markdown">
        <MarkdownHTML>{{ commentBody }}</MarkdownHTML>
      </div>
      <div class="flex-grow overflow-scroll text-sm border-l">
        <DiffBlock :content="comment.diff_hunk" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  comment: { type: Object, required: true },
  idx: { type: Number, required: true }
});

const commentBody = computed(() => {
 return props.comment.body.replace(/^\s*\d+\.\s*/, '');
});
</script>
