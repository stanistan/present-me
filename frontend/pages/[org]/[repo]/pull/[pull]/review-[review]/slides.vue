<template>
  <Review-PageChrome v-slot="data" height="h-[95vh]">  
    <SlideShow :model="data.data" />
  </Review-PageChrome>
</template>
