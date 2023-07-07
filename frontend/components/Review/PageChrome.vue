<template>
  <div>
    <TopBar>
      {{ $route.params.org }}/{{ $route.params.repo }}#{{ $route.params.pull }} 
      <div v-if="data" class="inline-block bg-slate-50 shadow-inner text-black px-2 py-1 rounded-sm text-xs">
        <ReviewLink :params="data.params" to="cards" /> | <ReviewLink :params="data.params" to="slides" />
      </div>
    </TopBar>
    <div class="relative" :class="height">
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
      <slot v-else :data="data" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  height: { type: String, default: "" }
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
