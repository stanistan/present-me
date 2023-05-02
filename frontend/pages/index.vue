<template>
  <div class="mx-auto">

    <div class="py-10 mb-10 bg-gray-100 border border-gray-200">
    <div class="text-5xl flex flex-col font-extrabold text-center py-3 h-[30vh]">
       <span class="flex-grow"></span>
      <GradientText class="flex-none">[pr]esent-me</GradientText>
    </div>

    <form @submit="submit" class="my-4 mx-auto text-lg max-w-3xl ">
      <div class="mx-2 flex flex-row
        rounded bg-white shadow-md
        p-2 gap-2
        border border-violet-100
        ">
        <input :disabled="formDisabled"
          name="search" type="text" :value="searchQuery"
          placeholder="$org/$repo/pull/$pull#pullrequestreview-$review"
          class="flex-grow px-4 font-mono
          focus:ring-none
          rounded overflow-hidden inline-block" />
        <button
          type="submit"
          :disabled="formDisabled"
          class="
            rounded p-4 px-6 text-lg font-bold bg-gradient-to-b from-purple-700 to-purple-800
            hover:from-purple-600 hover:to-purple-700
            active:from-purple-600 active:to-purple-300
            disabled:from-purple-200 disabled:to-purple-300 disabled:border-gray-100
            border border-gray-600 hover:border-gray-400
            text-white shadow-md">
          <span v-if="formDisabled" class="animate-pulse">
            ...
          </span>
          <span v-else>
            go
          </span>
        </button>
      </div>
    </form>
    </div>

    <div class="prose max-w-prose mx-auto gap-3 px-4">
      <div class="rounded-lg font-bold ring-1 my-5 ring-red-300 bg-red-100 p-3 text-center" v-if="errorMessage">
        Error: <span class="underline">{{ errorMessage }}</span>
      </div>

      <p class="inline-block font-bold">What</p>
      <p class="inline-block mb-4">
        <code>present-me</code> is an experiment to try to give the author of a Pull Request a better way to convey
        why a changeset looks the way that it does, and how the folks reading and reviewing it should approach it.
      </p>
      <p class="inline-block font-bold">How</p>
      <p class="mb-4">
        <code>present-me</code> uses a PR review's comments (and their respective diff) to create
        slideshow-like presentation, in the order that the comments are desired to appear, and only the
        diffs that are annotated with comments, leaving all other changes out of mind.
      </p>
      <p class="mb-2">
        These are all valid URLs to query for:
      </p>
      <ul class="list-disc ml-4 mb-4">
        <li v-for="(u, idx) in validURLs" class="text-sm">
          <strong>{{ u.why }}</strong> :: <br />
          <a href="#" @click="goTo(u.url)" class="text-xs underline">{{ u.url }}</a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'present-me'
});

const formDisabled = ref(false),
  errorMessage = ref(""),
  searchQuery = ref("");

async function goTo(url) {
  searchQuery.value = url;
  await executeSearch();
}

async function executeSearch() {
  formDisabled.value = true;
  errorMessage.value = "";
  setTimeout(async function() {
    const { data, error } = await useFetch('/api/search', {
      params: { search: searchQuery.value },
      server: false,
      initialCache: false,
      transform: v => JSON.parse(v)
    });

    if (error.value) {
      const errorData = JSON.parse(error.value.data)
      errorMessage.value = errorData.msg;
      formDisabled.value = false;
    } else {
      const params = data.value;
      await navigateTo(`${params.owner}/${params.repo}/pull/${params.number}/review-${params.review}`);
    }
  }, 2000);
}

async function submit(e) {
  e.preventDefault();
  await executeSearch();
}

const validURLs = [
  {
    why: 'Fully qualified Pull Request Review URL (the permalink from Github)',
    url: 'https://github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746',
  },
  {
    why: 'Dropping the Protocol (https is implicit)',
    url: 'github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746'
  },
  {
    why: 'Dropping the domain (https://github.com is implicit)',
    url: 'stanistan/invoice-proxy/pull/3#pullrequestreview-625362746'
  },
  {
    why: 'Dropping the URL fragment... will attempt to find the first Review by the PR author',
    url: 'stanistan/invoice-proxy/pull/3'
  }
];
</script>
