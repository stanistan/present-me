<!DOCTYPE html>
<html lang="en-US">
  <head>
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">

    {{ range .CSSFiles }}
    <link rel="stylesheet" href="{{ . }}">
    {{ end }}
  </head>
  <body>

    <div class="min-h-screen flex flex-col justify-content-center">
      <div class="flex-grow">
        {{ slot "body" }}
      </div>
      <div class="flex-none">

        <footer
          class="pt-3 pb-5 max-w-screen-2xl mx-auto font-mono text-xs flex flex-row"
        >
          <div class="flex-none px-3">
            <a href="{{ .Version.URL }}">({{ .Version.SHA }})</a>
          </div>
          <div class="flex-grow" />
          <div class="flex-none px-3">

            <a
              href="https://github.com/stanistan/present-me"
              class="underline underline-offset-8 font-bold">
              <em>present-me</em> <span>by stanistan</span>
            </a>

          </div>
        </footer>


      </div>
    </div>

    {{ range .JSFiles }}
    <script src="{{ . }}"></script>
    {{ end }}
  </body>
</html>
