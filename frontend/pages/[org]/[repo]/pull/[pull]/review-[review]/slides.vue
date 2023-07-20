<template>
  <ReviewPageChrome v-slot="data" height="h-[95vh]" name="slides">
    <SlideShow :model="data.data" />
  </ReviewPageChrome>
</template>
