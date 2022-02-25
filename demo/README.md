# filer-gateway demo scripts

The demo scripts are written in BASH.  They can be used to test the filer-gateway deployment and as a reference implementation.  The scripts require [`curl`](https://curl.haxx.se/) and [`jq`](https://stedolan.github.io/jq/).

To use the scripts, you need to know the API URL, API key, API user and API password of your filer-gateway deployment; and set them accordingly via the following environment variables:

- `API_URL`: the API base URL, default is set to `http://localhost:8080/v1`.
- `API_KEY`: the API key to be set to the header attribute `X-API-KEY` of the HTTP requests to the filer-gateway.
- `API_USER`: the API username for the HTTP Basic Authentication required by various POST/PATCH calls to the filer-gateway.
- `API_PASS`: THE API password for the HTTP Basic Authentication required by various POST/PATCH calls to the filer-gateway.

__An usage example can be found in the [`demo.sh`](demo.sh) file.__

The script [ping.sh](ping.sh) is written to check the API server's health by calling the `GET /ping` interface.  It also makes use of the OIDC shared token as an alternative authentication mechanism of the filer-gateway.  Therefore, instead of providing `API_KEY`, `API_USER` and `API_PASS`, one should provide the following env. variables:

- `API_AUTH_URL`: the authentication server URL 
- `API_CLIENT_ID`: the client id with scope "filer-gateway"
- `API_CLIENT_SECRET`: the client secret
