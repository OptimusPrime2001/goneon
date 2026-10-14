# Modules & Dependencies — Hướng dẫn chi tiết

Tài liệu tham khảo về **Go Modules** — hệ thống quản lý dependency chuẩn của Go
(theo `Roadmap.md` §8. Packages & Modules → phần *Modules*).

> Máy hiện tại: `go version go1.25.3 windows/amd64`

---

## 1. Khái niệm & Mục đích

**Module** là đơn vị **versioning và phân phối** cấp cao nhất trong Go: một tập hợp
package có chung một đường dẫn import, được khai báo bởi file `go.mod` ở thư mục gốc.

**Dependency** là module bên ngoài mà project của bạn cần để biên dịch hoặc chạy.

> 📌 **Module ≠ Package.** Module là "đóng gói để **version**" (thường = 1 repository);
> Package là "đóng gói để **import**" (1 thư mục có câu lệnh `package`). Một module chứa nhiều package.
> Chi tiết ở [Packages.md](<Packages.md>).

### Vì sao cần Modules?

Trước Modules (Go < 1.11) mọi thứ dựa vào `GOPATH` — không có khái niệm version:

| Vấn đề thời GOPATH | Modules giải quyết thế nào |
|---|---|
| Không có khái niệm version | Mỗi dependency có version tường minh trong `go.mod` |
| Lấy code từ `master`, dễ break | Chọn đúng tag `vX.Y.Z`, neo (pin) cố định |
| Hai project dùng 2 phiên bản khác nhau không dùng chung được | Module cache toàn cục, mỗi project tự có graph riêng |
| Không kiểm chứng nguồn code | `go.sum` + checksum DB chống supply-chain attack |

---

## 2. Cấu trúc & Quy tắc

### 2.1 Cấu trúc thư mục chuẩn

```text
goneon/
├── go.mod                 # bắt buộc – định danh module
├── go.sum                 # checksum – tự động sinh, commit vào repo
├── go.work                # (tùy chọn) workspace – KHÔNG publish
├── main.go
├── internal/              # chỉ import được trong module này
├── pkg/
└── cmd/
    └── server/
        └── main.go
```

### 2.2 Cú pháp `go.mod`

```go
module github.com/your-org/goneon     // đường dẫn import gốc (repository URL + path con)

go 1.22                               // phiên bản language tối thiểu

require (
    github.com/gin-gonic/gin v1.10.0       // direct dependency
    github.com/jackc/pgx/v5 v5.5.5
)

require (
    github.com/go-playground/validator/v10 v10.20.0 // indirect
)

replace github.com/old/lib => github.com/myfork/lib v1.2.3   // đổi nguồn

retract v1.0.1   // thu hồi version có lỗ hổng
```

| Directive | Ý nghĩa |
|---|---|
| `module` | Đường dẫn import gốc của chính module này |
| `go` | Phiên bản language tối thiểu (định `//go:build` behavior, feature gate) |
| `require` | Dependency bắt buộc. Dòng `// indirect` = dependency **của** dependency |
| `replace` | Trỏ import path sang nguồn khác (fork, path local) |
| `retract` | "Thu hồi" version đã phát hành (chứa bug/CVE) |
| `exclude` | Loại bỏ một version cụ thể khỏi graph |

### 2.3 Quy tắc version

| Quy tắc | Mô tả |
|---|---|
| **SemVer** | `vMAJOR.MINOR.PATCH` — bắt đầu bằng `v`, **bắt buộc** có `v` ở đầu |
| **MVS** *(Minimal Version Selection)* | Chọn phiên bản **nhỏ nhất** thỏa mãn toàn bộ graph |
| **Major ≥ 2** | Bắt buộc đổi import path: `github.com/x/lib/v2` |
| **Pseudo-version** | Chưa tag: `v0.0.0-20240101120000-abcdef123456` |
| **Pre-release** | `v1.2.0-rc.1`, `v1.2.0-beta.3` — không phải bản ổn định |

---

## 3. Hướng dẫn chi tiết & Ví dụ thực tế

### Bước 1 — Khởi tạo module

```bash
go mod init github.com/your-org/goneon
```

### Bước 2 — Thêm dependency

```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    r.Run(":8080")
}
```

```bash
go mod tidy   # tự quét source → thêm thiếu, xoá thừa
```

`go.mod` sau khi tidy:

```go
require github.com/gin-gonic/gin v1.10.0
```

### Bước 3 — Ép version cụ thể (override)

MVS chọn phiên bản **nhỏ nhất** — nếu dependency trung gian kẹp version cũ có bug,
bạn phải ép nó lên bằng một `require` trực tiếp:

```bash
# libB v1.2.0 (do libA yêu cầu) có bug → ép toàn bộ graph lên v1.4.0
go get github.com/x/libB@v1.4.0
```

```go
require github.com/x/libB v1.4.0   // require trực tiếp → thắng MVS
```

### Bước 4 — `replace` khi develop local

```go
// go.mod
replace github.com/your-org/gonice => ../gonice   // trỏ sang thư mục local
```

> ⚠️ **Không commit `replace` trỏ local** vào nhánh publish — người khác sẽ không build được.
> Dùng `go.work` thay thế cho mục đích dev liên repo.

### Bước 5 — `go.work` (multi-repo workspace)

```bash
go work init
go work use ./goneon
go work use ../gonice
```

```go
// go.work
go 1.22

use (
    ./goneon
    ../gonice
)
```

> 💡 `go.work` ảnh hưởng **riêng máy bạn**, không làm thay đổi `go.mod` của ai —
> lý tưởng để dev song song 2 repo mà không phải commit `replace`.

### Bước 6 — Kiểm tra dependency

```bash
go mod why github.com/jackc/pgpassfile   # tại sao package này nằm trong graph?
go mod graph                            # in dependency graph
go list -m all                          # toàn bộ module + version
```

---

## 4. Các lệnh quan trọng (CLI)

| Lệnh | Mục đích |
|---|---|
| `go mod init <path>` | Tạo `go.mod` mới |
| `go mod tidy` | Đồng bộ dependencies (thêm thiếu / xoá thừa) — **luôn chạy trước khi commit** |
| `go mod tidy -diff` | Chỉ in diff, không sửa (Go 1.23+) — dùng cho CI |
| `go get <pkg>@<ver>` | Thêm hoặc nâng cấp dependency |
| `go get -u ./...` | Upgrade toàn bộ dependency trong module |
| `go mod download` | Tải dependencies về module cache |
| `go mod verify` | Kiểm tra integrity qua `go.sum` |
| `go mod why <pkg>` | Giải thích vì sao package nằm trong graph |
| `go mod graph` | In dependency graph |
| `go mod vendor` | Sinh thư mục `vendor/` (đóng gói source) |
| `go mod edit -json` | Xuất `go.mod` dạng JSON |
| `go list -m all` | Liệt kê toàn bộ module + version |
| `go work init/use` | Quản lý workspace nhiều module |

> Danh sách lệnh đầy đủ: [Tooling/Commands.md](<../Tooling/Commands.md>)

---

## 5. Lỗi thường gặp & Best Practices

### ⚠️ Các bẫy phổ biến

| # | Bẫy | Cách tránh |
|---|---|---|
| 1 | Quên `go mod tidy` → CI fail `missing go.sum entry` | Thêm vào pre-commit hook |
| 2 | Nâng major version (`v1 → v2`) nhưng quên đổi import path | Đổi thành `.../lib/v2` ở **mọi** import |
| 3 | Commit `replace` trỏ local path | Dùng `go.work`; revert `replace` trước khi merge |
| 4 | `go.sum` bị sửa tay | **Không bao giờ** sửa tay — chỉ qua `go get` / `go mod tidy` |
| 5 | Depend vào pseudo-version (chưa tag) | Yêu cầu dependency có release tag ổn định |
| 6 | Upgrade ồ ạt `go get -u ./...` | Upgrade từng package, chạy test giữa chừng |
| 7 | Hiểu nhầm `v0.x` là "ổn định" | `v0.x` = API chưa cam kết; `v1.0.0` mới là stable |

### 💡 Best practices

* Commit cả `go.mod` **và** `go.sum`.
* CI chạy: `go mod tidy -diff && go mod verify`.
* Ưu tiên **ít dependency** — mỗi dependency là một gánh nặng bảo trì.
* Dùng `go mod vendor` nếu cần build offline/reproducible (đặt `GOFLAGS=-mod=vendor`).
* Upgrade dependency trong **PR riêng**, tách khỏi feature code.

---

**Bài liên quan**

* [Packages](<Packages.md>) — đơn vị tổ chức code
* [Package Import Rules](<Package Import Rules.md>) — luật import
* [Using Third Party Packages](<Using Third Party Packages.md>) — chọn & quản lý dependency ngoài
* [Publishing Modules](<Publishing Modules.md>) — phát hành module của riêng bạn
