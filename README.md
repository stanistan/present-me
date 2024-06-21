# present-me

> A tool to view/show a PR as a presentation / slide-show.

See https://prme.stanistan.com for a demo.

<img width="839" alt="image" src="https://github.com/stanistan/present-me/assets/66807/ba53a48f-a49e-4728-bef2-cfbfab318462">

## Why

Sometimes you make a PR that is too large for folks to read/parse top-to-bottom and it would be useful
to annotate it to desribe and show your team _how to read_ and parse the changeset.
What's important, what's superfluous, what's incidental to this change?

Present-me uses GitHub's Review Comments as a persistence layer (and write UI)
to generate something that can be, well, presented!

### How to use it

1. Create your PR as you would normally... [pr](https://github.com/stanistan/present-me/pull/56)
2. Start commenting on your own PR! But instead of individual comments, start a review!
3. Prefix your comments with a number to set the order that it'll show up
4. Grab the permalink of the review.
   <img width="680" alt="image" src="https://github.com/stanistan/present-me/assets/66807/89033c9c-6486-4da0-8cf9-d269443f0290">
5. Go to https://prme.stanistan.com/
6. Paste the permalink into the box
7. Hit GO! [rendered](https://prme.stanistan.com/stanistan/present-me/pull/56/review-1419621494)

<img width="797" alt="image" src="https://github.com/stanistan/present-me/assets/66807/1c0a6209-a135-4fde-aa12-d93c5316a4e8">

## Development

This project uses [hermit](https://cashapp.github.io/hermit/) for dependency management.

### Bootstrap

```sh
# for go packages
go mod download

# for tailwindcss
bun install
```

Make sure you've run this before anything else!

#### `.env`

This assumes you have a `.env` in the project directory.

```
PORT=8080
GH_APP_ID=0
GH_INSTALLATION_ID=0
LOG_OUTPUT=console
DISK_CACHE_ENABLED=true
DISK_CACHE_BASE_PATH=data
DISK_CACHE_MAX_SIZE=10000
```

N.B. Not having the ID, InstallationID, and PK file will use the public GH API and be subject
to those rate limits. Enabling the `DISK_CACHE` helps alleviate those. You can check out
the cached files in `data/` to see what the responses look like.

Full options:

```
# go run ./cmd/veun --help
Usage: present-me --log-output="json" --gh-app_id=INT-64 --gh-installation_id=INT-64 <command>

Flags:
  -h, --help                              Show context-sensitive help.
      --port="8080"                       ($PORT)
      --hostname="localhost"              ($HOSTNAME)
      --server-read-timeout=5s
      --server-write-timeout=10s
      --debug                             ($DEBUG)
      --environment="dev"                 ($ENV)
      --log-output="json"                 ($LOG_OUTPUT)
      --disk-cache-enabled                ($DISK_CACHE_ENABLED)
      --disk-cache-base-path=STRING       ($DISK_CACHE_BASE_PATH)
      --disk-cache-cache-max-size=1024    ($DISK_CACHE_MAX_SIZE_KB)
      --gh-app_id=INT-64                  ($GH_APP_ID)
      --gh-installation_id=INT-64         ($GH_INSTALLATION_ID)
      --gh-pk-file=STRING                 ($GH_PK_FILE)

Commands:
  version --log-output="json" --gh-app_id=INT-64 --gh-installation_id=INT-64

  serve --log-output="json" --gh-app_id=INT-64 --gh-installation_id=INT-64
```

#### Development

```bash
dev
```

- <http://localhost:8080>
