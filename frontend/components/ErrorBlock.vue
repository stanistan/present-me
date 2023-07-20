<template>
  <div class="mx-auto max-w-3xl py-10">
    <div class="bg-orange-100 px-2 pb-2">
      <div class="text-xs text-center underline py-4">
        <slot name="title" />
      </div>
      <div class="bg-white p-4 text-center border border-orange-200 rounded">
        <slot />
      </div>
    </div>
  </div>
</template>
