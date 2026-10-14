# Packages — Hướng dẫn chi tiết

Tài liệu tham khảo về **Package system** trong Go
(theo `Roadmap.md` §8. Packages & Modules → phần *Package System*).

> Đọc trước: [Module and Dependencies](<Module and Dependencies.md>) để phân biệt **Module** và **Package**.

---

## 1. Khái niệm & Mục đích

**Package** là đơn vị **tổ chức code** nhỏ nhất có thể import: một thư mục chứa các file
Go cùng một câu lệnh `package <tên>`.

```go
package geometry   // mọi file trong thư mục này phải khai báo cùng một package
```

### Mục đích

* **Encapsulation** — chỉ export những gì cần (identifier viết hoa).
* **Đóng gói logic** thành module tái sử dụng, kiểm thử được.
* **Điều khiển phạm vi tầm nhìn** qua `internal/` và tên file `_test.go`.
* **Tái sử dụng** — package dùng chung giữa nhiều binary trong cùng module.

> 📌 **Module vs Package**

| | Package | Module |
|---|---|---|
| Định nghĩa | thư mục có `package foo` | cây thư mục có `go.mod` |
| Purpose | đơn vị **tổ chức code**, `import` | đơn vị **versioning** & dependency |
| Một module | gồm **nhiều** package | — |
| Một package | thuộc **duy nhất** một module | — |

`github.com/gin-gonic/gin` là **package**; module chứa nó cũng là
`github.com/gin-gonic/gin` (vì `go.mod` nằm ở root repo).

---

## 2. Cấu trúc & Quy tắc

### 2.1 Quy tắc đặt tên Package

| Quy tắc | ✅ Đúng | ❌ Sai |
|---|---|---|
| Ngắn, viết thường, không underscore | `httputil`, `jsonutil` | `http_utils`, `JSONUtil` |
| Một từ nếu có thể | `strings`, `sync` | `string_helpers` |
| Không trùng package chuẩn | `myhttp` | `http` (gây lẫn lộn khi import) |
| Tên thư mục ≈ tên package | `net/http` → `package http` | thư mục `utils` nhưng `package helpers` |
| Package `main` dùng cho binary | `package main` | dùng `main` cho thư viện |

> 📌 **Trừu tượng hóa:** người dùng import theo **đường dẫn thư mục**,
> nhưng gọi hàm theo **tên package**. Tên package ngắn gọn → code đọc hay hơn:
> `http.Get(...)`, `strings.ToUpper(...)`, `json.Marshal(...)`.

### 2.2 Cấu trúc thư mục khuyến nghị

```text
github.com/your-org/goneon
│
├── cmd/                          # binary entrypoints – mỗi thư mục 1 main package
│   └── server/
│       └── main.go               # package main
│
├── internal/                     # ❌ KHÔNG import được từ ngoài module
│   └── auth/
│       ├── auth.go               # package auth
│       └── auth_test.go
│
├── pkg/                          # (tùy chọn) API công khai, import bởi repo khác
│   └── client/
│       └── client.go             # package client
│
├── api/                          # OpenAPI / protobuf definitions
├── configs/
└── go.mod
```

| Thư mục | Vai trò | Import từ ngoài module |
|---|---|---|
| `cmd/<bin>/` | Entry point của binary (`package main`) | ❌ (là program) |
| `internal/` | Code riêng của module, API không ổn định | ❌ **bị cấm** |
| `pkg/` | API công khai, cam kết ổn định | ✅ |
| thư mục khác | Package theo domain | ✅ |

> 💡 `pkg/` **không phải chuẩn bắt buộc** của Go — nó chỉ là tín hiệu "đây là API public".
> Có team để trống root bằng các package domain (`/user`, `/order`) thay vì `pkg/`.

### 2.3 Quy tắc file & visibility

| Identifier | Ký hiệu | Phạm vi |
|---|---|---|
| Exported | `Func`, `Type`, `Const` | Mọi nơi import được package |
| Unexported | `func`, `type`, `const` | Chỉ trong cùng package |
| Kiểm thử nội bộ | `foo_test.go` → `package foo` | Truy cập được unexported |
| Kiểm thử bên ngoài | `foo_test.go` → `package foo_test` | Chỉ dùng API public (khuyến nghị) |

```go
package telemetry

type Config struct {          // exported struct
    Endpoint string          // exported field
    sample   float64         // unexported field – không truy cập từ package khác
}

func NewClient() {}           // exported
func parse() {}               // unexported
```

### 2.4 Package `init()`

```go
func init() {
    // chạy TRƯỚC main(), theo thứ tự import (depth-first)
    registerDriver("mysql", mysqlDriver)
}
```

> ⚠️ Dùng `init()` tối giản — khó test, khó debug, dễ tạo side-effect ẩn.
> Ưu tiên hàm `New(...)` tường minh.

---

## 3. Hướng dẫn chi tiết & Ví dụ thực tế

### 3.1 Cấu trúc module hoàn chỉnh

```text
goneon/
├── go.mod
├── cmd/
│   └── server/main.go
├── internal/
│   ├── config/
│   │   ├── config.go        // package config
│   │   └── config_test.go   // package config (internal test)
│   └── repository/
│       └── user.go          // package repository
└── pkg/
    └── telemetry/
        └── telemetry.go     // package telemetry
```

### 3.2 `cmd/server/main.go` — package `main`

```go
package main

import (
    "log"

    "github.com/your-org/goneon/internal/config"
    "github.com/your-org/goneon/internal/server"
)

func main() {
    cfg, err := config.Load("configs/app.yaml")
    if err != nil {
        log.Fatalf("load config: %v", err)
    }
    if err := server.Run(cfg); err != nil {
        log.Fatal(err)
    }
}
```

> 📌 `package main` **không import được** từ package khác:
> compiler báo `is a program, not an importable package`.
> Toàn bộ logic nghiệp vụ phải nằm ở package con, `main` chỉ **wiring**.

### 3.3 Package API public

```go
// pkg/telemetry/telemetry.go
package telemetry

import "fmt"

// Config cấu hình cho telemetry client.
type Config struct {
    Endpoint string
    Sample   float64
}

// NewClient tạo client mới. Mọi tên exported đều phải có doc comment.
func NewClient(cfg Config) (*Client, error) {
    if cfg.Endpoint == "" {
        return nil, fmt.Errorf("telemetry: endpoint is required")
    }
    return &Client{cfg: cfg}, nil
}

type Client struct{ cfg Config }   // field unexported → chỉ thao tác qua method
```

### 3.4 Test cùng package (truy cập unexported)

```go
// internal/config/config_test.go
package config

import "testing"

func TestParse_Defaults(t *testing.T) {
    got := parse("")          // hàm unexported – chỉ gọi được trong package
    if got.Timeout == 0 {
        t.Fatal("default timeout must not be zero")
    }
}
```

### 3.5 Test API public (khuyến nghị)

```go
// pkg/telemetry/telemetry_test.go
package telemetry_test      // package test KHÁC → chỉ dùng API exported

import (
    "testing"

    "github.com/your-org/goneon/pkg/telemetry"
)

func TestNewClient_RequiresEndpoint(t *testing.T) {
    if _, err := telemetry.NewClient(telemetry.Config{}); err == nil {
        t.Fatal("expected error for empty endpoint")
    }
}
```

> 💡 Test `package foo_test` mô phỏng **người dùng thật** của library —
> bắt được lỗi quên export, sai godoc, API đổi ngầm.

---

## 4. Các lệnh quan trọng (CLI)

| Lệnh | Mục đích |
|---|---|
| `go build ./...` | Build toàn bộ package |
| `go test ./...` | Chạy test toàn bộ |
| `go test ./internal/... -cover` | Test + coverage |
| `go test ./pkg/... -race -cover` | Test có race detector |
| `go vet ./...` | Static analysis, bắt lỗi phổ biến |
| `go doc <pkg>` | Xem documentation của package |
| `go doc <pkg>.<Sym>` | Xem documentation của symbol |
| `gofmt -l .` | Liệt kê file chưa format |
| `go list ./...` | Liệt kê package trong module |
| `golangci-lint run` | Lint tổng hợp (team) |

```bash
# Xem nhanh API một package
go doc github.com/your-org/goneon/pkg/telemetry
go doc github.com/your-org/goneon/pkg/telemetry.NewClient
```

---

## 5. Lỗi thường gặp & Best Practices

### ⚠️ Các bẫy phổ biến

| # | Bẫy | Hệ quả | Cách tránh |
|---|---|---|---|
| 1 | 2 thư mục 2 package name khác nhau | Lẫn lộn khi đọc code | Tên package = tên thư mục |
| 2 | Đặt tên package `utils`, `common`, `helper` | Không rõ responsibility, thành "rác" | Tách theo domain: `httputil`, `timeutil` |
| 3 | Package `main` chứa logic nghiệp vụ | Không test được, không tái sử dụng | Logic ở package riêng, `main` chỉ wiring |
| 4 | Package quá to (hàng nghìn LOC) | Khó maintain, dễ circular | "Small, focused packages" |
| 5 | Circular import giữa 2 package | **Biên dịch fail** `import cycle not allowed` | Tách package thứ ba hoặc dùng interface |
| 6 | Export thừa | API surface phình, khó thay đổi | Bắt đầu unexported, export khi có người dùng |
| 7 | Đổi tên file package name | `expected package X` | Giữ đồng bộ `package X` ↔ tên thư mục |
| 8 | Quên doc comment cho tên exported | Code kẹt `golint`, khó dùng | Viết comment mô tả **đầu tiên** |

### 💡 Best practices

* Viết **doc comment** cho mọi identifier exported — dùng `golint`/`revive`.
* Test cùng package (`package foo`) khi cần test nội bộ;
  test `package foo_test` để bảo vệ API public.
* Đặt **package `main` mỏng** — chỉ parse flag, đọc config, gọi vào package app.
* Ưu tiên **"small, focused packages"** — triết lý Go.
* Chạy `go vet ./...` + `golangci-lint run` trong CI.

---

**Bài liên quan**

* [Package Import Rules](<Package Import Rules.md>) — luật import & chống cycle
* [Module and Dependencies](<Module and Dependencies.md>) — version & dependency
