# Package Import Rules — Hướng dẫn chi tiết

Tài liệu tham khảo về **luật import package** trong Go
(theo `Roadmap.md` §8. Packages & Modules → phần *Package System*).

> Đọc trước: [Packages](<Packages.md>)

---

## 1. Khái niệm & Mục đích

**Import** là cơ chế tham chiếu package khác để dùng symbol của nó. Go có **luật import
chặt chẽ**: compiler phân tích toàn bộ graph import để đảm bảo không có cycle,
không có unused import.

### Mục đích

* Minh bạch nguồn gốc mọi symbol đang dùng.
* Ngăn **circular dependency**.
* Kiểm soát tầm nhìn qua `internal/`.
* `gofmt` / `goimports` tự động chuẩn hoá thứ tự import.

```go
import "github.com/gin-gonic/gin"     // khai báo → dùng gin.Context trong file
```

---

## 2. Cấu trúc & Quy tắc

### 2.1 Cú pháp

```go
import "github.com/gin-gonic/gin"                       // đơn

import (
    "fmt"                                                // 1. stdlib
    "net/http"

    "github.com/gin-gonic/gin"                          // 2. bên thứ ba
    "github.com/jackc/pgx/v5"

    "github.com/your-org/goneon/internal/config"        // 3. module của bạn
)
```

### 2.2 Thứ tự import (Standard Grouping)

| Nhóm | Nội dung | Cách phân tách |
|---|---|---|
| 1 | **stdlib** (`fmt`, `net/http`, `context`) | dòng trắng |
| 2 | **External modules** (đường dẫn bắt đầu bằng domain) | dòng trắng |
| 3 | **Local / module nội bộ** | dòng trắng |

> 📌 `goimports` tự sắp xếp và nhóm theo chuẩn này. **Không sắp tay** — để tool quyết định.
> CI nên check: `test -z "$(gofmt -l .)" && test -z "$(goimports -l .)"`.

### 2.3 Các dạng import

| Dạng | Cú pháp | Ý nghĩa | Khuyến nghị |
|---|---|---|---|
| Thường | `"strings"` | Gọi `strings.ToUpper` | ✅ |
| Named alias | `g "github.com/gin-gonic/gin"` | Tránh trùng tên package | ⚠️ Khi cần |
| Blank import | `_ "net/http/pprof"` | Chạy side-effect của `init()` | ✅ Hợp lệ |
| Dot import | `. "fmt"` | Gọi thẳng `Printf` | ❌ **Cấm trong team** |

```go
// ✅ Blank import – kích hoạt side-effect (register driver, mount handler...)
import (
    _ "github.com/lib/pq"        // register DB driver vào database/sql
    _ "net/http/pprof"           // mount /debug/pprof lên DefaultServeMux
)

// ⚠️ Alias – CHỈ khi trùng tên THẬT SỰ
import (
    "encoding/json"              // json.Marshal
    sjson "github.com/tidwall/sjson"   // alias để phân biệt
)

// ❌ Dot import – không bao giờ dùng trong codebase chung
import (
    . "fmt"                      // Printf(...) → không rõ Printf từ đâu ra
)
```

### 2.4 Quy tắc `internal/`

```text
your-org/goneon/internal/auth     ← CHỈ package trong your-org/goneon import được
```

* Import từ **ngoài module** → lỗi biên dịch:
  `use of internal package ... not allowed`
* Dùng cho code **không muốn trở thành public API** (chi tiết ở [Packages](<Packages.md>)).

### 2.5 Các quy tắc cấm (build fail)

| Quy tắc | Thông báo lỗi |
|---|---|
| ❌ Circular import `A → B → A` | `import cycle not allowed` |
| ❌ Import rồi không dùng | `imported and not used` |
| ❌ Import package `main` | `is a program, not an importable package` |
| ❌ Import vượt `internal/` | `use of internal package not allowed` |
| ❌ Relative import (`./foo`) | Không tồn tại trong Modules |
| ❌ Import path sai major version | `no required module provides package ...` |

---

## 3. Hướng dẫn chi tiết & Ví dụ thực tế

### 3.1 File import chuẩn

```go
// internal/handler/user.go
package handler

import (
    "context"
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/your-org/goneon/internal/service"
    "github.com/your-org/goneon/pkg/telemetry"
)

func GetUser(c *gin.Context) {
    ctx := context.Background()
    u, err := service.FindUser(ctx, c.Param("id"))
    if err != nil {
        _ = telemetry.Record(err)
        c.JSON(http.StatusNotFound, gin.H{"error": errors.New("not found").Error()})
        return
    }
    c.JSON(http.StatusOK, u)
}
```

### 3.2 Xử lý trùng tên package

```go
import (
    "encoding/json"                       // json.Marshal

    // ⚠️ KHÔNG import cùng package 2 lần (duplicate import)
    // "github.com/tidwall/sjson"         // ❌ alias phải THAY cho path thường
    sjson "github.com/tidwall/sjson"      // ✅ alias rõ nghĩa khi cần phân biệt
)

func demo() {
    b, _ := json.Marshal(map[string]int{"a": 1})
    out, _ := sjson.SetBytes(b, "b", 2)
    _, _ = b, out
}
```

> 💡 **Nguyên tắc:** alias **chỉ khi xung đột tên thật sự**.
> Alias tràn lan làm code khó truy vết symbol.

### 3.3 Tránh circular import bằng interface

```go
// ❌ SAI: repository import service, service import repository → CYCLE
//
//   service/user.go   → import repo
//   repo/user.go      → import service

// ✅ ĐÍNH: định nghĩa interface ở package CONSUMER (dependency inversion)

// internal/service/user.go
package service

import "context"

type UserRepo interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type Service struct {
    repo UserRepo        // service KHÔNG import repository
}

func New(r UserRepo) *Service { return &Service{repo: r} }

// internal/repository/user.go
package repository

import (
    "context"

    "github.com/your-org/goneon/internal/service"
)

// Repo triển khai service.UserRepo (structural typing – không cần implements)
func (r *Repo) FindByID(ctx context.Context, id string) (*service.User, error) {
    // ...
    return nil, nil
}

// cmd/server/main.go – wiring ở tầng ngoài
svc := service.New(repository.NewPool(dsn))
```

```text
Trước (CYCLE):          service ⇄ repository
Sau (DI):               repository → service      (Repo implement interface)
                        main → cả hai              (wiring ở tầng ngoài)
```

### 3.4 Import pattern cho test

```go
// Pattern thường gặp: table-driven test import package cần test
import "testing"

var tests = []struct{
    name string
    in   string
    want string
}{
    {"empty", "", ""},
    {"basic", "a", "A"},
}

func TestUpper(t *testing.T) {
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := upper(tt.in); got != tt.want {
                t.Errorf("upper(%q) = %q, want %q", tt.in, got, tt.want)
            }
        })
    }
}
```

---

## 4. Các lệnh quan trọng (CLI)

| Lệnh | Mục đích |
|---|---|
| `goimports -w .` | Format + tự thêm/bỏ import |
| `goimports -l .` | Liệt kê file import chưa chuẩn (CI) |
| `gofmt -w .` | Chỉ format, không sửa import |
| `gofmt -l .` | Liệt kê file chưa format |
| `go mod tidy` | Loại bỏ import không còn dùng / thêm import thiếu |
| `go build ./...` | Kiểm tra import hợp lệ |
| `go list -deps ./...` | In toàn bộ dependency của package |
| `go list -f '{{.ImportPath}}: {{.Imports}}' ./...` | Import của từng package |
| `golangci-lint run` | Lint tổng hợp (gồm `gci`, `misspell`…) |

```bash
# Cài goimports (một lần)
go install golang.org/x/tools/cmd/goimports@latest

# Check nhanh trước khi commit
gofmt -l . && goimports -l . && go vet ./...
```

---

## 5. Lỗi thường gặp & Best Practices

### ⚠️ Các bẫy phổ biến

| # | Lỗi | Thông báo | Cách sửa |
|---|---|---|---|
| 1 | Import rồi không dùng | `imported and not used` | Xoá import hoặc chạy `goimports` |
| 2 | Cycle giữa 2 package | `import cycle not allowed` | Tách interface / package thứ ba |
| 3 | Import `internal` từ ngoài | `use of internal package not allowed` | Chuyển sang `pkg/` nếu cần public |
| 4 | Import path sai version major | `no required module provides package ...` | Thêm `/v2`, `/v3` vào path |
| 5 | Dot import | Code khó truy vết symbol | **Cấm tuyệt đối** |
| 6 | Sắp import thủ công | Tranh luận trong review, merge conflict | Để `goimports` quyết định |
| 7 | Import vòng qua `main` | `is a program, not an importable package` | Tách logic ra package con |
| 8 | Gọi init side-effect nhưng quên blank import | Symbol undefined / init không chạy | Dùng `_ "pkg"` tường minh |

### 💡 Best practices

* CI: `test -z "$(gofmt -l .)" && test -z "$(goimports -l .)"` → fail nếu lệch.
* Alias **chỉ khi xung đột tên thật sự**.
* Không export chỉ để "đỡ viết import nội bộ" — dùng `internal/`.
* Import nhóm rõ ràng giúp review diff nhanh hơn nhiều.
* Khi gặp `import cycle not allowed`: **đừng** phá `internal/` để lách —
  hãy suy nghĩ lại ranh giới package.

---

**Bài liên quan**

* [Packages](<Packages.md>) — cấu trúc & visibility
* [Using Third Party Packages](<Using Third Party Packages.md>) — quản lý dependency ngoài
