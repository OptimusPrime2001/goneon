# Using Third Party Packages — Hướng dẫn chi tiết

Tài liệu tham khảo về **sử dụng dependency bên thứ ba** trong Go
(theo `Roadmap.md` §8. Packages & Modules → *Dependency management*).

> Đọc trước: [Module and Dependencies](<Module and Dependencies.md>)

---

## 1. Khái niệm & Mục đích

**Third party package** là module do cộng đồng / tổ chức khác phát hành, lấy qua
**Go Module Proxy** (`proxy.golang.org`) hoặc trực tiếp từ VCS (Git).

**Mục đích:** tái sử dụng năng lực đã có (HTTP framework, ORM, logger…) thay vì viết lại —
nhưng **mỗi dependency là một rủi ro**: bảo trì, bảo mật, licensing, transitive weight.

> 💡 **Nguyên tắc vàng**
> *"A little copying is better than a little dependency."* — Go Blog.
> Viết 10 dòng code nội bộ đôi khi rẻ hơn nhiều so với kéo cả một dependency tree.

---

## 2. Cấu trúc & Quy tắc

### 2.1 Checklist chọn dependency

| Tiêu chí | Câu hỏi cần trả lời |
|---|---|
| **Activity** | Commit gần nhất < 6 tháng? Maintainer phản hồi issue? |
| **Adoption** | Có team lớn dùng production không? Stars, downloads? |
| **API stability** | Đã `v1.x` ổn định? Breaking change thường xuyên? |
| **License** | MIT / Apache-2.0 / BSD? ⚠️ GPL với product closed-source |
| **Transitive weight** | Kéo theo bao nhiêu `// indirect`? |
| **Security** | Có CVE history? `govulncheck` sạch? |
| **Alternatives** | Stdlib đã làm được chưa? Package `golang.org/x/...` đã đủ chưa? |

### 2.2 Quy tắc quản lý version

| Quy tắc | Chi tiết |
|---|---|
| **Pin version** | Luôn commit `require ... vX.Y.Z` cụ thể |
| **SemVer** | `v1.x.y` upgrade an toàn; `v2+` cần đổi import path |
| **MVS** | Version thực dùng = **nhỏ nhất** thỏa mọi ràng buộc trong graph |
| **Review upgrade** | Đọc `CHANGELOG` trước khi `go get -u` |
| **Direct vs indirect** | Càng ít **direct** dependency càng tốt |

### 2.3 Biến môi trường liên quan

```bash
export GOPROXY=https://proxy.golang.org,direct   # proxy mặc định
export GOPRIVATE=github.com/your-org/*            # bỏ qua proxy + sum DB cho repo nội bộ
export GONOPROXY=github.com/your-org/*            # chỉ bỏ qua proxy
export GONOSUMDB=github.com/your-org/*            # chỉ bỏ qua sum DB
export GOFLAGS=-mod=vendor                        # build luôn dùng vendor/
export GOMODCACHE=$HOME/go/pkg/mod                # module cache toàn cục
```

| Biến | Vai trò |
|---|---|
| `GOPROXY` | Nơi tải module (`proxy.golang.org`, `direct`, `off`) |
| `GOPRIVATE` | Shorthand set **cả** `GONOPROXY` + `GONOSUMDB` |
| `GONOSUMDB` | Không đối chiếu checksum với sum DB |
| `GOFLAGS` | Cờ mặc định cho mọi lệnh `go` |

> ⚠️ **Private repo:** phải set `GOPRIVATE=github.com/your-org/*`
> nếu không `go get` sẽ fail vì proxy không truy cập được repo riêng tư
> và sum DB không verify được.

---

## 3. Hướng dẫn chi tiết & Ví dụ thực tế

### Bước 1 — Thêm dependency

```bash
go get github.com/jackc/pgx/v5@v5.5.5
```

```go
// internal/repository/user.go
package repository

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
    return pgxpool.New(ctx, dsn)
}

var _ = pgx.ErrNoRows
```

### Bước 2 — Đồng bộ `go.mod` / `go.sum`

```bash
go mod tidy
```

```go
// go.mod
require github.com/jackc/pgx/v5 v5.5.5

require (
    github.com/jackc/pgpassfile v1.0.0            // indirect
    github.com/rogpeppe/go-internal v1.12.0        // indirect
    golang.org/x/crypto v0.23.0                    // indirect
)
```

> 📌 `// indirect` = dependency của dependency — bạn không import trực tiếp.
> `go mod tidy` tự động thêm/xoá các dòng này.

### Bước 3 — Audit & truy vết

```bash
go mod why github.com/jackc/pgpassfile    # tại sao package này xuất hiện?
go list -m -u all                         # có version mới không?
go list -m -json github.com/jackc/pgx/v5  # chi tiết 1 module
govulncheck ./...                          # quét lỗ hổng đã biết
```

### Bước 4 — Vendor (tuỳ chọn, build reproducible)

```bash
go mod vendor
go build -mod=vendor ./...
```

```text
vendor/
├── github.com/jackc/pgx/v5/...
├── modules.txt            # manifest – commit vào repo
└── ...
```

> 💡 Dùng `vendor/` khi: build offline, CI cần reproducible tuyệt đối,
> hoặc nội bộ không tin tưởng proxy. Chi phí: repo phình lên.

### Bước 5 — `replace` cho fork bảo mật khẩn cấp

```go
// go.mod – fork vá CVE, gỡ khi upstream release
replace github.com/vulnerable/lib => github.com/your-org/lib v1.0.1-security.0
```

> ⚠️ `replace` là **giải pháp tạm**. Phải có ticket theo dõi để gỡ
> khi upstream vá xong — quên sẽ giữ fork vĩnh viễn.

### Quy trình nâng cấp dependency an toàn

```bash
# 1. Tạo nhánh riêng
git checkout -b chore/upgrade-pgx

# 2. Đọc CHANGELOG trước
#    https://github.com/jackc/pgx/releases

# 3. Upgrade từng package
go get github.com/jackc/pgx/v5@v5.6.0
go mod tidy

# 4. Chạy test
go test ./... -race

# 5. Audit
govulncheck ./...
go mod verify

# 6. Merge – tách khỏi feature code
```

---

## 4. Các lệnh quan trọng (CLI)

| Lệnh | Mục đích |
|---|---|
| `go get <pkg>@<ver>` | Thêm / đổi version dependency |
| `go get -u <pkg>` | Upgrade lên minor/patch mới nhất |
| `go get -u ./...` | Upgrade toàn bộ (cẩn thận) |
| `go get -t ./...` | Bao gồm dependency của test |
| `go mod tidy` | Đồng bộ `go.mod`/`go.sum` với source |
| `go mod download` | Tải dependencies về cache |
| `go mod verify` | Kiểm tra checksum |
| `go mod why <pkg>` | Giải thích nguồn gốc dependency |
| `go mod graph` | In dependency graph |
| `go list -m -u all` | Liệt kê + version mới nhất |
| `go list -m -json <mod>` | Chi tiết 1 module |
| `govulncheck ./...` | Quét lỗ hổng bảo mật |
| `go mod vendor` | Sinh thư mục `vendor/` |

```bash
# Cài công cụ audit (một lần)
go install golang.org/x/vuln/cmd/govulncheck@latest
```

---

## 5. Lỗi thường gặp & Best Practices

### ⚠️ Các bẫy phổ biến

| # | Bẫy | Rủi ro | Cách tránh |
|---|---|---|---|
| 1 | "Đi chợ" — thêm dependency cho việc 5 dòng code làm được | Phình graph, tăng attack surface | Viết trước, refactor sau |
| 2 | Upgrade major không đọc CHANGELOG | Breaking API âm thầm | Đọc CHANGELOG + chạy test |
| 3 | Dependency bị abandon (2 năm không commit) | Lỗ hổng không được vá | Fork hoặc thay thế |
| 4 | Dùng `latest` / branch `main` | Không reproducible | Luôn pin tag `vX.Y.Z` |
| 5 | Bỏ qua license | Vấn đề pháp lý | Check LICENSE trước khi merge |
| 6 | Không chạy `govulncheck` | Để lọt CVE | Chạy định kỳ + trong CI |
| 7 | `replace` tạm bị quên khi merge | Build người khác fail | Ticket tracking + checklist review |
| 8 | Trả `// indirect` bằng tay | Lệch với graph thật | Chỉ để tool sửa `go.mod` |
| 9 | Upgrade ồ ạt `go get -u ./...` | Không biết cái nào break | Upgrade từng package |

### 💡 Best practices

* **Dependency review trong PR:** mỗi dependency mới phải được justify
  (thêm dòng vào PR template: *"Tại sao cần package này?"*).
* CI: `go mod tidy -diff && go mod verify && govulncheck ./...`.
* Ưu tiên **`golang.org/x/...`** trước khi cần util
  (`x/sync`, `x/time`, `x/net`) — chính chủ Go, bảo trì tốt.
* Ghép upgrade dependency thành **PR riêng**, tách khỏi feature code.
* Giữ dependency **direct** tối thiểu — kiểm tra định kỳ bằng `go mod why`.

---

**Bài liên quan**

* [Module and Dependencies](<Module and Dependencies.md>) — cấu trúc `go.mod`, MVS
* [Publishing Modules](<Publishing Modules.md>) — khi tự bạn phát hành module
