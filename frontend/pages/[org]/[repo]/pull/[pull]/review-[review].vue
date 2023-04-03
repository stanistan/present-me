<template>
    <div class="mx-auto">
        <div class="bg-black text-white text-center">
            {{ $route.params.org }} / {{ $route.params.repo }} / {{ $route.params.pull }} / {{ $route.params.review }}
        </div>

        <div class="whitespace-pre-line">
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
  transform: v => JSON.parse(v),
});
</script>
