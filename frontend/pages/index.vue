<template>
  <div class="mx-auto">
    <div class="py-10 mb-10 bg-gray-100 border border-gray-200">
      <div class="text-5xl flex flex-col font-extrabold text-center py-3 h-[30vh]">
        <span class="flex-grow" />
        <GradientText class="flex-none">
          [pr]esent-me
        </GradientText>
      </div>

      <SearchBox
        v-model="query" :error-message="errorMessage"
        :disabled="searchDisabled" @submit="search"
      />
    </div>

    <div class="prose max-w-prose mx-auto px-4">
      <p class="inline-block font-bold">
        What
      </p>
      <p class="inline-block mb-4">
        <code>present-me</code> is an experiment to try to give the author of a Pull Request a better way to convey
        why a changeset looks the way that it does, and how the folks reading and reviewing it should approach it.
      </p>
      <p class="inline-block font-bold">
        How
      </p>
      <p class="mb-4">
        <code>present-me</code> uses a PR review's comments (and their respective diff) to create
        slideshow-like presentation, in the order that the comments are desired to appear, and only the
        diffs that are annotated with comments, leaving all other changes out of mind.
      </p>
      <p class="mb-2">
        These are all valid URLs to query for:
      </p>
      <ul class="list-disc ml-4 mb-4">
        <li
          v-for="u in validURLs" :key="u.url" 
          class="text-sm"
        >
          <strong>{{ u.why }}</strong> :: <br>
          <a
            href="#" class="text-xs underline"
            @click="goTo(u.url)"
          >{{ u.url }}</a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'present-me'
});

const query = ref(''),
  errorMessage = ref(''),
  searchDisabled = ref(false);

function searchLoading() {
  searchDisabled.value = true;
  errorMessage.value = '';
}

function searchError(msg) {
  searchDisabled.value = false;
  errorMessage.value = msg;
}

async function search() {
  searchLoading();

  setTimeout(async function() {
    const { data, error } = await useFetch('/api/search', {
      params: { search: query.value },
      server: false,
      initialCache: false,
      transform: v => JSON.parse(v)
    });

    if (error.value) {
      searchError(JSON.parse(error.value.data).msg);
    } else {
      const params = data.value;
      await navigateTo(`${params.owner}/${params.repo}/pull/${params.number}/review-${params.review}`);
    }
  }, 1000);
}

async function goTo(url) {
  query.value = url;
  await search();
}

const validURLs = [
  {
    why: 'Fully qualified Pull Request Review URL (the permalink from Github)',
    url: 'https://github.com/stanistan/present-me/pull/56#pullrequestreview-1419621494',
  },
  {
    why: 'Dropping the Protocol (https is implicit)',
    url: 'github.com/stanistan/present-me/pull/56#pullrequestreview-1419621494'
  },
  {
    why: 'Dropping the domain (https://github.com is implicit)',
    url: 'stanistan/present-me/pull/56#pullrequestreview-1419621494'
  },
  {
    why: 'Dropping the URL fragment... will attempt to find the first Review by the PR author',
    url: 'stanistan/present-me/pull/56'
  }
];
</script>
