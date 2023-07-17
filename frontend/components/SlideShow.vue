<template>
  <div class="flex flex-col h-full" @keyup.left="left">
    <div class="flex-grow" />
    <div class="flex-0 max-w-[2200px] mx-auto">
      <div v-if="current==0">
        <div class="text-6xl font-extrabold text-center"> 
          <span>(#{{ model.pr.number }})</span>&nbsp;
          <GradientText>{{ model.pr.title }}</GradientText>
        </div>
        <div class="mx-auto mt-8">
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

      <div v-for="(c, i) in model.comments" :key="i">
        <SlideCard
          v-if="i+2==current" :comment="c"
          :idx="i+1"
        />
      </div>

      <div v-if="current===model.comments.length+2" class="text-center font-bold">
        FIN
      </div>
    </div>
    <div class="flex-grow" />
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  model: { type: Object, required: true }
});

function onKeyUp(e) {
  if (e.defaultPrevented) {
    return;
  }

  e.preventDefault();

  const totalSlides = props.model.comments.length + 3;
  let next = current.value;

  switch (e.key) {
    case "ArrowLeft": 
      next = (next - 1) % totalSlides;
      break;
    case "ArrowRight":
    case "Space":
      next = (next + 1) % totalSlides;
      break;
  }


  if (next < 0) {
    next = totalSlides - 1; 
  }

  current.value = next;
}

onMounted(() => {
  window.addEventListener('keyup', onKeyUp);
});

onUnmounted(() => {
  window.removeEventListener('keyup', onKeyUp);
});

const current = ref(0);

</script>