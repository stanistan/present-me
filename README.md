# present-me

## Local Development

### Prerequisites

#### Bootstrap

Make sure you've run `prmectl bootstrap`!

#### `.env`

This assumes you have a `.env` in your project root.

```
GH_APP_ID=
GH_INSTALLATION_ID=
PORT=8080
GH_PK_FILE=path-to-cert.pem
```

### `prmectl dev`

This will start both the go server at port `8080` and
the nuxt client at port `3000`.

You can interact with both via port `8080` since the go
server will proxy directly to nuxt.

```bash
prmectl dev
```
