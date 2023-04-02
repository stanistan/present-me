<script setup lang="ts">

const route = useRoute();
const { pending, data } = await useFetch('/api/review', {
  params: route.params,
  server: false,
  initialCache: false,
  transform: v => JSON.parse(v),
});

</script>

<template>
    <div>
        <div class="bg-black text-white text-center">
            {{ $route.params.org }} / {{ $route.params.repo }} / {{ $route.params.pull }} / {{ $route.params.review }}
        </div>

        <div class="whitespace-pre-line">

          <pre> {{ pending ? '...' :  data.Params.Owner }}</pre>
        </div>
    </div>
</template>
