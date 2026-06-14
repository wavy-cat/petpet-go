# petpet-go

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/wavy-cat/petpet-go?style=for-the-badge&logo=go&logoColor=white&labelColor=1A222E&color=242B36)
![GitHub License](https://img.shields.io/github/license/wavy-cat/petpet-go?style=for-the-badge&labelColor=1A222E&color=242B36)
![GitHub repo size](https://img.shields.io/github/repo-size/wavy-cat/petpet-go?style=for-the-badge&logo=github&logoColor=white&labelColor=1A222E&color=242B36&cacheSeconds=0)

A web service for generating petpet GIFs based on a Discord user's avatar, written in Go.

## Usage

### Discord avatar

<kbd>GET</kbd> `/ds/{user_id}.gif`

#### Path parameters

| Name        | Type      | Description           |
|-------------|-----------|-----------------------|
| `{user_id}` | Snowflake | The Discord user's ID |             

#### Query parameters

| Name       | Default | Type             | Description                                        |
|------------|---------|------------------|----------------------------------------------------|
| `delay`    | `4`     | Unsigned Integer | GIF speed. Bigger is slower                        |
| `no-cache` | `false` | Boolean          | Whether to disable caching (Cache-Control headers) |

### Custom upload

<kbd>POST</kbd> `/custom`

Send a `multipart/form-data` request with a file field named `image` containing a PNG, JPEG, WebP or AVIF.

The upload is limited to 5MB and a maximum of 1 MP (you can reconfigure it).

## Getting started

* Static binaries available as [releases](https://github.com/wavy-cat/petpet-go/releases)
* [GHCR](https://github.com/wavy-cat/petpet-go/pkgs/container/petpet-go) or `docker pull ghcr.io/wavy-cat/petpet-go` (see [Compose](compose.yml) example)

> [!NOTE]
> The Discord bot sends preview requests from the `us-east1` GCP region (South Carolina, US). To reduce network latency, choose server locations close to it. If using a Cloudflare proxy, the nearest Cloudflare location is ATL (Atlanta, US).

## Configuration

Currently, config parameters can be specified either in the `config.yml` file or via environment variables.

Outbound HTTP proxy settings are handled by Go's standard environment variables: `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY`.

Look at the [configuration reference](config.sample.yml) with comments (including environment variable names).

## PetPet in other languages

* **Python**: [nakidai/petthecord](https://github.com/nakidai/petthecord)
* **Rust**: [mitsuki-kagamin/petpet-api](https://github.com/mitsuki-kagamin/petpet-api)
* **C**: [nakidai/cptc](https://github.com/nakidai/cptc)
