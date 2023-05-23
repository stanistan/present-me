<template>
  <div>
    <TopBar>
      {{ $route.params.org }}/{{ $route.params.repo }}/pull/{{ $route.params.pull }}/review-{{ $route.params.review }}
    </TopBar>
    <div v-if="pending" class="flex flex-col items-stretch">
      <div class="animate-pulse mx-auto text-center text-4xl pt-10 font-bold ">
        <GradientText>Loading...</GradientText>
      </div>
    </div>
    <div v-else-if="error" class="mx-auto max-w-3xl py-10">
      <div class="bg-orange-100 px-2 pb-2">
        <div class="text-xs text-center underline py-4">
          {{ error }}
        </div>
        <div class="bg-white p-4 text-center border border-orange-200 rounded">
          <code>{{ error.data }}</code>
        </div>
      </div>
    </div>
    <div v-else class="mx-auto max-w-screen-2xl">
      SLIDES WILL GO HERE
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'present-me',
});

const route = useRoute();
const { pending, data, error } = await useFetch('/api/review', {
  lazy: true,
  params: route.params,
  server: false,
  initialCache: false,
  transform: v => JSON.parse(v),
});
</script>
