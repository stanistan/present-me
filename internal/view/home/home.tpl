<div class="mx-auto">
    <div class="pt-10 mb-10 bg-indigo-100 border-b border-gray-200">
      <div class="text-5xl flex flex-col font-extrabold text-center pt-3">
        [pr]esent-me
      </div>
      <div class="prose max-w-prose mx-auto px-6 text-center text-xs py-2">
        <p>
          <code>present-me</code> uses the comments on a Pull Request to generate a slideshow.
        </p>
      </div>
      <div class="pt-5 prose mx-auto text-center">
        <p class="underline underline-offset-8">
          Present a Pull Request.
        </p>
      </div>
      <div class="bg-indigo-900 mt-3">
        <form action="/" method="GET">
        <div class="flex items-center gap-1 font-mono justify-center">
          <div class="py-4">
            <span class="text-gray-400">//github.com/</span>
          </div>
          <div class="py-4">
            <input type="text" name="owner" placeholder="owner" value="{{ .Owner }}" class="
              border-t border-t-gray-900
              shadow-lg
              w-28
              text-center p-1 rounded text-xs
            ">
          </div>
          <div class="py-4">
            <span class="text-gray-400">/</span>
          </div>
          <div class="py-4">
            <input type="text" name="repo" placeholder="repo" value="{{ .Repo }}" class="
              border-t border-t-gray-900
              shadow-lg
              w-28
              text-center p-1 rounded text-xs
            ">
          </div>
          <div class="py-4">
            <span class="text-gray-400">/pull/</span>
          </div>
          <div class="py-4">
            <input type="text" name="pull" placeholder="pr" value="{{ .Pull }}" class="
              border-t border-t-gray-900
              shadow-lg
              text-center p-1 rounded text-xs w-12
            ">
          </div>
          <div>
            <button class="
              px-4
              pb-[6px] pt-[4px]
              ml-1
              bg-purple-700 rounded-md
              hover:bg-purple-600
              active:bg-purple-800
              text-white text-sm font-bold font-sans
              shadow-xs hover:shadow-md
              border-b border-b-gray-400 border-t border-t-purple-500"
              type="submit"
              >
              pr-me
            </button>
          </div>
        </div>
        </form>
      </div>
    {{ slot "results" }}
    </div>
    <div class="prose max-w-prose mx-auto px-4">
      <p class="mb-4">
        There are a few kinds of PRs out there in the wild. There are the PRs that
        are trivial to review, they are small, self-contained, etc.
      </p>
      <p>
        There are others though:
      </p>
        <ul class="list-disc ml-4">
        <li class="prose py-1">
          You might be updating a function signature in a library and have to update
          the call site across hundreds of files. There's too much <em>noise</em> to
          find the parts that are most important.
        </li>
        <li class="prose py-1">
          You might be refactoring something more involved and want to tell a story
          with your PR, or it's just an idea and you want to walk through it.
        </li>
        <li class="prose py-1">
          You might have left a lot of comments annotating your PR but how do you
          make sure that folks read them all in the right order? They're going to
          look at the diff and read it <em>top to bottom</em> anyway.
        </li>
        </ul>
    </div>

  </div>
</div>
