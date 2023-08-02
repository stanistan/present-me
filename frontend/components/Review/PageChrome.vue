<template>
  <div>
    <TopBar>
      {{ $route.params.org }}/{{ $route.params.repo }}#{{ $route.params.pull }}
      <div
        v-if="data"
        class="inline-block bg-slate-50 shadow-inner text-black px-2 py-1 rounded-sm text-xs"
      >
        <ReviewLink :params="data.meta.params" to="cards" :current="name" /> |
        <ReviewLink :params="data.meta.params" to="slides" :current="name" />
      </div>
      <template v-if="name == 'slides'" #right>
        <button
          class="text-xs px-2 text-violet-300 hover:text-pink-300"
          @click="requestFullscreen"
        >
          :play:
        </button>
      </template>
    </TopBar>
    <div class="relative" :class="height">
      <div v-if="pending" class="flex flex-col items-stretch">
        <div class="animate-pulse mx-auto text-center text-4xl pt-10 font-bold">
          <GradientText>Loading...</GradientText>
        </div>
      </div>
      <ErrorBlock v-else-if="error">
        <template #title>
          {{ error }}
        </template>
        <code>{{ error.data }}</code>
      </ErrorBlock>
      <div v-else ref="content" class="h-full">
        <slot :data="data!!" name="default" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Review } from "../../src/Review";
defineProps<{
  height?: string;
  name: string;
}>();

const route = useRoute();
const { pending, data, error } = await useFetch<Review>("/api/review", {
  lazy: true,
  params: route.params,
  server: false,
});

const content = ref<Element>()!!;
const requestFullscreen = () => {
  content.value!!.requestFullscreen();
};
</script>
