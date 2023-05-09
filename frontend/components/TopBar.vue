<template>
  <div class="bg-gradient-to-b from-gray-800 to-black text-white font-mono text-sm py-2 shadow">
    <div class="flex flex-row">
      <div class="flex-none pl-3">
        <NuxtLink href="/" class="font-bold hover:text-pink-300">/</NuxtLink>
      </div>
      <div class="flex-grow"></div>
      <div class="flex-none">
        <slot />
      </div>
      <div class="flex-grow"></div>
    </div>
  </div>
</template>
