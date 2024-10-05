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

| Name                    | Default   | Example            | Description                                                                                                                                        |
|-------------------------|-----------|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| `ADDRESS`               | `:80`     | `127.0.0.1:443`    | The address (including port) where the server will run.                                                                                            |
| `PORT`                  | `80`      | `443`              | The port where the server will run. Used if `ADDRESS` is not specified.                                                                            |
| `SHUTDOWN_TIMEOUT`      | `5`       | `10`               | Time in seconds for correct server shutdown                                                                                                        |
| `BOT_TOKEN`             | null      | `your_token`       | Token from Discord bot. Used for authorization in the Discord API.                                                                                 |
| `CACHE_STORAGE`         | null      | `memory`           | The storage type used for caching images. Possible values: `memory` (LRU cache in RAM), `fs` (file system). If not specified, caching is disabled. |
| `CACHE_MEMORY_CAPACITY` | `100`     | `1984`             | The memory storage capacity (maximum number of items). This option is used when `CACHE_STORAGE` is set to `memory`.                                |
| `CACHE_FS_PATH`         | `./cache` | `/mnt/s3fs/petpet` | The path to the directory used for file system-based cache storage.                                                                                |

> [!NOTE]
> `BOT_TOKEN` is an optional variable. If you don't specify it, the server will receive user avatars not directly from
> Discord, but through avatar.cdev.shop.

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
