<template>
    <div class="mx-auto">
        <div class="bg-gradient-to-b from-gray-800 to-black text-white font-mono text-sm text-center py-2 shadow">
          {{ $route.params.org }}/{{ $route.params.repo }}/pull/{{ $route.params.pull }}/review-{{ $route.params.review }}
        </div>
        <div>
          <div v-if="pending">
            ... loading ...
          </div>
          <div v-else>
            <Review :model="data" />
          </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const route = useRoute();
const { pending, data } = await useFetch('/api/review', {
  params: route.params,
  server: false,
  initialCache: false,
  transform: v => JSON.parse(v)
});
</script>
