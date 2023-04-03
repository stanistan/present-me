# present-me

## Local development...

### Prerequisites

#### Bootstrap

```sh
prmectl bootstrap
```

Make sure you've run  before anything else!

#### `server/.env`

This assumes you have a `.env` in the server directory.

```
GH_APP_ID=
GH_INSTALLATION_ID=
GH_PK_FILE=path-to-cert.pem
```

### Development

```bash
prmectl dev
```

This will start both the go server at port `8080` and
the nuxt client at port `3000`.

You can interact with both via port `8080` since the go
server will proxy directly to nuxt.

- <http://localhost:8080>

## Building an image for production...

