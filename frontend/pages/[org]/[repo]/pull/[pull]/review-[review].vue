<script setup lang="ts">
const route = useRoute();
const data = await useFetch('/api/review', {
    params: route.params,
    server: false,
    initialCache: false
});
</script>
<template>
    <div>
        <div>
            {{ $route.params.org }} / {{ $route.params.repo }} / {{ $route.params.pull }} / {{ $route.params.review }}
        </div>
        <div>
        <pre> {{ data.data }}</pre>
        </div>
    </div>
</template>
