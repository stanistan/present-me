<template>
  <div class="mx-auto max-w-4xl">
    <div class="text-5xl font-extrabold text-center py-3">
      <span class="bg-clip-text text-transparent bg-gradient-to-r underline from-pink-600 to-violet-900">
        [pr]esent-me
      </span>
    </div>

    <form @submit="submit" class="mt-4 text-lg">
      <div class="mx-auto flex flex-row
        rounded bg-white shadow-md
        p-2 gap-2
        border border-violet-100
        ">
        <input ref="searchBox" :disabled="formDisabled"
          name="search" type="text"
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
            border border-gray-600 hover:border-gray-400
            text-white shadow-md">
          go
        </button>
      </div>
      <div class="rounded-lg font-bold ring-1 mt-5 ring-red-300 bg-red-100 p-3 text-center" v-if="errorMessage">
        Error: <span class="underline">{{ errorMessage }}</span>
      </div>
    </form>

    <div class="prose mt-4 max-w-prose mx-auto gap-3 px-4">
      <p class="inline-block mb-4">
        (pr)esent-me is an experiment to try to give the author of a Pull Request a better way to convey
        why a changeset looks the way that it does, and how the folks reading and reviewing it should approach it.
      </p>
      <p class="inline-block font-bold mb-4">How it works</p>
      <p class="mb-4">
        <code>present-me</code> uses a PR review's comments (and their respective diff) to create a single
        "post", or "slides."
      </p>
      <p class="">
        These are all valid URLs to query for:
        <ul class="list-disc">
          <li class="text-sm">
            <strong>Fully qualified Pull Request Review URL (the permalink from Github)</strong> :: <br />
            <code class="text-xs">https://github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746</code>
          </li>
          <li class="text-sm">
            <strong>Dropping the Protocl (https is implicit)</strong> :: <br />
            <code class="text-xs">github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746</code>
          </li>
          <li class="text-sm">
            <strong>Dropping the domain (https://github.com is implicit)</strong> :: <br />
            <code class="text-xs">stanistan/invoice-proxy/pull/3#pullrequestreview-625362746</code>
          </li>
          <li class="text-sm">
            <strong>Dropping the URL fragment... This will attempt to find the first PR review by the author</strong> :: <br />
            <code class="text-xs">stanistan/invoice-proxy/pull/3</code>
          </li>
        </ul>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">

const searchBox = ref(),
      formDisabled = ref(false),
      errorMessage = ref("");

onMounted(() => {
  console.log(searchBox.value.focus());
})

async function submit(e) {
  e.preventDefault();
  formDisabled.value = true;
  const { data, error } = await useFetch('/api/search', {
    params: { search: searchBox.value.value },
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
}
</script>
