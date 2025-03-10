# petpet-go

A web service for generating petpet GIFs based on a Discord user's avatar.

## Getting started

* Docker

```bash
docker run ghcr.io/wavy-cat/petpet-go
```

* Go

```bash
go run github.com/wavy-cat/petpet-go/cmd/app
```

## Environment

| Env Name                | YAML Name                | Type                                    | Default   | Example        | Description                                                                                                                                        |
|-------------------------|--------------------------|-----------------------------------------|-----------|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| `HOST`                  | `server.host`            | string                                  | null      | `127.0.0.1`    | The address where the server will run.                                                                                                             |
| `PORT`                  | `server.port`            | uint16                                  | `3000`    | `8080`         | The port where the server will run. Used if `ADDRESS` is not specified.                                                                            |
| `SHUTDOWN_TIMEOUT`      | `server.shutdownTimeout` | [uint](https://pkg.go.dev/builtin#uint) | `5000`    | `10`           | Time in milliseconds for correct server shutdown                                                                                                   |
| `BOT_TOKEN`             | `discord.botToken`       | string                                  | null      | Your bot token | Bot token from Discord. Used for authorization in the Discord API.                                                                                 |
| `CACHE_STORAGE`         | `cache.storage`          | One of `memory`, `fs`                   | null      | `memory`       | The storage type used for caching images. Possible values: `memory` (LRU cache in RAM), `fs` (file system). If not specified, caching is disabled. |
| `CACHE_MEMORY_CAPACITY` | `cache.memoryCapacity`   | [uint](https://pkg.go.dev/builtin#uint) | `100`     | `1984`         | The memory storage capacity (maximum number of items). This option is used when `CACHE_STORAGE` is set to `memory`.                                |
| `CACHE_FS_PATH`         | `cache.fsPath`           | string                                  | `./cache` | `/mnt/petpet`  | The path to the directory used for file system-based cache storage.                                                                                |

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
