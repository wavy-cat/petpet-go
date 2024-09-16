# petpet-go

Веб-сервер для генерации petpet гифок на основе аватарки пользователя Discord.

## Getting started

* Go

```bash
go run github.com/wavy-cat/petpet-go/cmd/petpet-go
```

## Usage

<kbd>GET</kbd> `/ds/{user_id}?delay=5`

`{user_id}` - ID пользователя в Discord.

`?delay` - скорость GIF (по-умолчанию 5).