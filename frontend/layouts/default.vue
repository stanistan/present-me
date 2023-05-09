<template>
  <div class="min-h-screen flex flex-col">
    <div class="flex-grow">
      <slot />
    </div>
    <div class="flex-none">
      <footer class="pt-3 pb-5 text-center text-sm font-mono text-xs">
        <NuxtLink 
          href="https://github.com/stanistan/present-me" 
          class="underline underline-offset-8 font-bold"
        >
          <em>present-me</em> by stanistan
        </NuxtLink> 
      </footer>
    </div>
  </div>
</template>
