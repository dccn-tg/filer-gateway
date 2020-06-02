# Demo scripts

The demo scripts in this directory make use of [`curl`](https://curl.haxx.se/) and [`jq`](https://stedolan.github.io/jq/).

To use the scripts, you need to know the API URL, API key, API user and API password of your filer-gateway deployment; and set them accordingly via the following environment variables:

- `API_URL`: the API base URL, default is set to `http://localhost:8080/v1`.
- `API_KEY`: the API key to be set to the header attribute `X-API-KEY` of the HTTP requests to the filer-gateway.
- `API_USER`: the API username for the HTTP Basic Authentication required by various POST/PATCH calls to the filer-gateway.
- `API_PASS`: THE API password for the HTTP Basic Authentication required by various POST/PATCH calls to the filer-gateway.
