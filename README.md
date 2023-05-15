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
prmectl bootstrap
```

Make sure you've run this before anything else!

#### `server/.env`

This assumes you have a `.env` in the server directory.

```
PORT=8080
GH_APP_ID=
GH_INSTALLATION_ID=
GH_PK_FILE=path-to-cert.pem
```

#### Development

```bash
prmectl dev
```

This will start both the go server at port `8080` and
the nuxt client at port `3000`.

You can interact with both via port `8080` since the go
server will proxy directly to nuxt.

- <http://localhost:8080>

## Testing a production binary

```bash
prmectl local-prod
```

This will generate the static output from `nuxt` and run the server
in production mode (not a proxy to the nuxt dev-server).
