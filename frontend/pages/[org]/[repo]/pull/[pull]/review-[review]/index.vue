<template>
  <Review-PageChrome v-slot="data">
    <div class="mx-auto max-w-screen-2xl">
      <ReviewPageContent :model="data.data" />
    </div>
  </Review-PageChrome>
</template>
