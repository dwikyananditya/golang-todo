# todo

Service todo sederhana — Go, Fiber v3, GORM, SQLite, validasi pakai go-playground/validator.

## Jalankan

```sh
go run ./cmd/todo
```

Listen di `:8080`, tabel `todo.db` (SQLite) otomatis di-migrate.

## API

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/` | Ambil semua todo |
| GET | `/:id` | Ambil todo by id |
| POST | `/` | Bikin todo |
| PATCH | `/:id` | Update todo |
| DELETE | `/:id` | Hapus todo |

`title` dan `description` wajib diisi.

## Contoh

```sh
curl -X POST localhost:8080/ -H 'Content-Type: application/json' \
  -d '{"title":"beli susu","description":"","is_done":false}'

curl localhost:8080/
```
