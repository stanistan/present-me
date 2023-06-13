<template>
  <div class="min-h-screen flex flex-col justify-content-center">
    <div class="flex-grow">
      <slot />
    </div>
    <div class="flex-none">
      <footer class="pt-3 pb-5 text-center font-mono text-xs">
        <div>
          <NuxtLink href="https://github.com/stanistan/present-me" class="underline underline-offset-8 font-bold">
            <em>present-me</em> by stanistan
          </NuxtLink> 
        </div>
        <div class="my-4">
          <NuxtLink :href="'https://github.com/stanistan/present-me/commit/' + v.rev">
            ({{ v. rev }})
          </NuxtLink>
        </div>
      </footer>
    </div>
  </div>
</template>

<script type="ts" setup>
import * as v from '../version.json';
</script>
