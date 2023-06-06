<template>
  <TopBar>
    {{ $route.params.org }}/{{ $route.params.repo }}/pull/{{ $route.params.pull }}/review-{{ $route.params.review }}
  </TopBar>
  <div class="relative h-[96vh]">
    <div v-if="pending">
      <div class="flex flex-col items-stretch">
        <div class="animate-pulse mx-auto text-center text-4xl pt-10 font-bold ">
          <GradientText>Loading...</GradientText>
        </div>
      </div>
    </div>
    <div v-else-if="error">
      <div class="mx-auto max-w-3xl py-10">
        <div class="bg-orange-100 px-2 pb-2">
          <div class="text-xs text-center underline py-4">
            {{ error }}
          </div>
          <div class="bg-white p-4 text-center border border-orange-200 rounded">
            <code>{{ error.data }}</code>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="absolute top-0 left-0 right-0 bottom-0">


    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute();
const { pending, data, error } = await useFetch('/api/review', {
  lazy: true,
  params: route.params,
  server: false,
  initialCache: false,
  transform: v => JSON.parse(v),
});
</script>
