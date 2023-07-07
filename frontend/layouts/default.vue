<template>
  <div class="min-h-screen flex flex-col justify-content-center">
    <div class="flex-grow">
      <slot />
    </div>
    <div class="flex-none">
      <footer class="pt-3 pb-5 max-w-screen-2xl mx-auto font-mono text-xs flex flex-row">
        <div class="flex-none px-3">
          <NuxtLink :href="'https://github.com/stanistan/present-me/commit/' + v.rev">
            ({{ revision }})
          </NuxtLink>
        </div>
        <div class="flex-grow" />
        <div class="flex-none px-3">
          <NuxtLink href="https://github.com/stanistan/present-me" class="underline underline-offset-8 font-bold">
            <em>present-me</em> by stanistan
          </NuxtLink> 
        </div>
      </footer>
    </div>
  </div>
</template>

<script type="ts" setup>
import * as v from '../version.json';
const revision = v.rev.slice(0, 7);
</script>
