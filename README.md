# petpet-go

Веб-сервер для генерации petpet гифок на основе аватарки пользователя Discord.

## Getting started

* Go

```bash
go run github.com/wavy-cat/petpet-go/cmd/petpet-go
```

## Environment

| Name      | Default | Descripton                                                                     |
|-----------|---------|--------------------------------------------------------------------------------|
| `ADDRESS` | `:80`   | Адрес (включая порт), на котором будет работать сервер                         |
| `PORT`    | `80`    | Порт на котором будет работать сервер. Используется, если `ADDRESS` не задано. |

## Usage

<kbd>GET</kbd> `/ds/{user_id}?delay=5&no-cache=false`

Параметры:

`{user_id}` - ID пользователя в Discord.

`?delay` (int) - скорость GIF. По-умолчанию `5`.

`?no-cache` (bool) - отключить ли кэширование (заголовки Cache-Control). По-умолчанию `false`.