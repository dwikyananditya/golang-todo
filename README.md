# todo

> **Catatan:** versi lama (HTTP/Fiber) ada di branch [`main`](https://github.com/dwikyananditya/golang-todo/tree/main).

Service todo sederhana pakai Go — gRPC + GORM + SQLite, dengan validasi request via protovalidate.

## Prasyarat

- Go 1.27+
- [buf](https://buf.build) dengan `protoc-gen-go` dan `protoc-gen-go-grpc` di PATH
- grpcurl (untuk testing)

## Setup

```sh
buf generate   # regenerate pb dari todo/v1/todo.proto
go build ./...
```

## Jalankan

```sh
go run ./cmd/todo
```

Listen di `:8080`, bikin/migrasi `todo.db` (SQLite) di direktori kerja.
Shutdown graceful saat SIGINT/SIGTERM.

## API

| RPC | Deskripsi |
|-----|-----------|
| `CreateTodo` | Bikin todo (title, description wajib) |
| `GetTodo` | Ambil todo by id |
| `ListTodos` | Ambil semua todo dengan sorting, filter, pagination |
| `UpdateTodo` | Update by id |
| `DeleteTodo` | Hapus by id |

Parameter `ListTodos`:

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `order` | enum | `SORT_ORDER_ASC` (default) atau `SORT_ORDER_DESC` |
| `is_done` | optional bool | Filter status selesai; kalau tidak diisi = tanpa filter |
| `limit` | int32 | Ukuran halaman, default 50, max 100 |
| `offset` | int32 | Skip N baris |

## Coba

```sh
grpcurl -plaintext -d '{"title":"beli susu","description":"","is_done":false}' \
  localhost:8080 todo.v1.TodoService/CreateTodo

grpcurl -plaintext -d '{"order":"SORT_ORDER_DESC","is_done":false,"limit":10}' \
  localhost:8080 todo.v1.TodoService/ListTodos
```

Reflection aktif, jadi tool seperti grpcurl dan Yaak bisa menemukan API tanpa file proto.
