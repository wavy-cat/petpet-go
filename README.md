# petpet-go

A web service for generating petpet GIFs based on a Discord user's avatar.

## Getting started

* Docker

```bash
docker run ghcr.io/wavy-cat/petpet-go
```

* Go

```bash
go run github.com/wavy-cat/petpet-go/cmd/petpet-go
```

## Environment

| Name               | Default | Example         | Description                                                             |
|--------------------|---------|-----------------|-------------------------------------------------------------------------|
| `ADDRESS`          | `:80`   | `127.0.0.1:443` | The address (including port) where the server will run.                 |
| `PORT`             | `80`    | `443`           | The port where the server will run. Used if `ADDRESS` is not specified. |
| `SHUTDOWN_TIMEOUT` | `5`     | `10`            | Time in seconds for correct server shutdown                             |
| `BOT_TOKEN`        | None    | `your_token`    | Token from Discord bot. Used for authorization in the Discord API.      |

> [!NOTE]
> `BOT_TOKEN` is an optional variable. If you don't specify it, the server will receive user avatars not directly from Discord, but through avatar.cdev.shop.

## Usage

<kbd>GET</kbd> `/ds/{user_id}?delay=5&no-cache=false`

### Path parameters

| Name        | Type      | Description            |
|-------------|-----------|------------------------|
| `{user_id}` | Snowflake | The Discord user's ID. |             

### Query parameters

| Name       | Default | Type             | Description                                         |
|------------|---------|------------------|-----------------------------------------------------|
| `delay`    | `5`     | Unsigned Integer | GIF speed.                                          |
| `no-cache` | `false` | Boolean          | Whether to disable caching (Cache-Control headers). |  

## PetPet in Other Languages

* **Python** - [nakidai/petthecord](https://github.com/nakidai/petthecord)
* **Rust** - [messengernew/petpet-api](https://github.com/messengernew/petpet-api)
