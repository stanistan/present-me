<template>
  <div class="flex flex-col h-full">
    <div class="flex-grow"></div>
    <div class="flex-0"> 

      <div v-if="current==0">
        <div class="text-6xl font-extrabold text-center"> 
          <span>(#{{ model.pr.number }})</span>&nbsp;
          <GradientText>{{ model.pr.title }}</GradientText>
        </div>
        <div class="mx-auto w-1/2 mt-8">
          <Review-MetadataList :model="model" />
        </div>
      </div>

      <div v-if="current===1">
        <ComponentCard> 
        <template #title>
          <div class="text-xl font-extrabold">
            <span>(#{{ model.pr.number }})</span>&nbsp;
            <GradientText>{{ model.pr.title }}</GradientText>
          </div>
        </template>
        <template #body>
          <div class="px-4 py-4">
            <BodyMarkdown :model="model" />
          </div>
        </template>
        </ComponentCard>
      </div>

      <div v-for="(c, i) in model.comments">
        <SlideCard :comment="c" :idx="i+1" v-if="i==current-2"/>
      </div>
    </div>
    <div class="flex-grow"></div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  model: { type: Object, required: true }
});

const current = ref(1);
</script>
