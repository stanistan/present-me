<script setup lang="ts">
const data = await useFetch('/api/ping', {
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
            {{ data.data }}
        </div>
    </div>
</template>
