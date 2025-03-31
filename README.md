# petpet-go

A web service for generating petpet GIFs (and APNG now) based on a Discord user's avatar.

## Usage

<kbd>GET</kbd> `/ds/{user_id}.gif`

### Path parameters

| Name        | Type      | Description           |
|-------------|-----------|-----------------------|
| `{user_id}` | Snowflake | The Discord user's ID |             

### Query parameters

| Name       | Default | Type             | Description                                        |
|------------|---------|------------------|----------------------------------------------------|
| `delay`    | `4`     | Unsigned Integer | GIF speed. Bigger is slower                        |
| `no-cache` | `false` | Boolean          | Whether to disable caching (Cache-Control headers) |

### Formats

* `.gif`
* `.apng`

## Getting started

* Docker (Container Registry)

```bash
docker run ghcr.io/wavy-cat/petpet-go
```

* Docker (Local)

```bash
docker build . -t ghcr.io/wavy-cat/petpet-go
docker run ghcr.io/wavy-cat/petpet-go
```

* Go

```bash
go run github.com/wavy-cat/petpet-go/cmd/app
```

## Configuration

Currently, config parameters can be specified either in the `config.yml` file or via environment variables.

Look at the [sample config](config.sample.yml) with comments (including environment variable names)

## PetPet in other languages

* **Python** - [nakidai/petthecord](https://github.com/nakidai/petthecord)
* **Rust** - [messengernew/petpet-api](https://github.com/messengernew/petpet-api)
