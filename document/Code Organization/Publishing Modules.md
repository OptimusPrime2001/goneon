# Publishing Modules — Hướng dẫn chi tiết

Tài liệu tham khảo về **phát hành (publish) một Go module** để team/cộng đồng `go get` được
(theo `Roadmap.md` §8. Packages & Modules).

> Đọc trước: [Module and Dependencies](<Module and Dependencies.md>)

---

## 1. Khái niệm & Mục đích

**Publishing module** = để người khác `go get` được module của bạn:
cấu trúc repo đúng chuẩn, gắn tag SemVer, và để **Go Module Proxy** lan toả.

### Mục đích

* Chia sẻ library nội bộ giữa các team / repo.
* Phân phối binary & CLI tool qua `go install`.
* Thiết lập versioning, giao tiếp breaking change rõ ràng.

> 📌 **Không cần "đăng ký" trung tâm.**
> Publish = đẩy code lên Git + gắn tag SemVer. Proxy tự lấy qua VCS.

```text
git tag v1.0.0  ──▶  git push origin v1.0.0  ──▶  proxy.golang.org  ──▶  go get
```

---

## 2. Cấu trúc & Quy tắc

### 2.1 Yêu cầu bắt buộc

| # | Yêu cầu | Chi tiết |
|---|---|---|
| 1 | **Repo public** (hoặc cấu hình proxy cho private) | GitHub / GitLab / Gitea… |
| 2 | **`go.mod` ở root** | `module github.com/you/lib` = URL repo |
| 3 | **Tag SemVer** | `v1.0.0`, `v1.2.3`, `v2.0.0` |
| 4 | **Tag bắt đầu bằng `v`** | `v1.0.0` ✅ · `1.0.0` ❌ |
| 5 | **Không còn `replace` local** | Sẽ làm người dùng build fail |
| 6 | **README + LICENSE** | MIT / Apache-2.0 / BSD phổ biến |

### 2.2 Semantic Import Versioning (bắt buộc)

| Version | `module` khai báo | Import path |
|---|---|---|
| `v0.x.x` | `github.com/you/lib` | `github.com/you/lib` |
| `v1.x.x` | `github.com/you/lib` | `github.com/you/lib` |
| `v2.x.x` | `github.com/you/lib/v2` | `github.com/you/lib/v2` |
| `vN.x.x` (N≥2) | `github.com/you/lib/vN` | `github.com/you/lib/vN` |

> ⚠️ **Major ≥ 2 bắt buộc đổi import path.**
> Đây là lý do Go có thể tồn tại song song `lib` (v1) và `lib/v2` (v2)
> trong cùng một dependency graph.

### 2.3 Quy tắc đặt version (SemVer)

| Loại thay đổi | Bump | Ví dụ |
|---|---|---|
| Fix bug, không đổi API | `PATCH` | `v1.0.0 → v1.0.1` |
| Thêm API tương thích | `MINOR` | `v1.0.1 → v1.1.0` |
| Xoá/sửa API, đổi behavior | `MAJOR` | `v1.9.0 → v2.0.0` |
| Chưa ổn định, đang dev | `v0.x` | `v0.3.0` (không cam kết stability) |
| Thử nghiệm | pre-release | `v1.2.0-rc.1`, `v1.2.0-beta.3` |

---

## 3. Hướng dẫn chi tiết & Ví dụ thực tế

### Bước 1 — Cấu trúc repo publish-ready

```text
github.com/your-org/gonice/
├── go.mod                     # module github.com/your-org/gonice
├── README.md                  # quickstart + ví dụ dùng
├── LICENSE                    # MIT / Apache-2.0
├── CHANGELOG.md               # ghi mọi release
├── .golangci.yml
├── client/                    # package public
│   └── client.go
├── internal/                  # KHÔNG publish API
│   └── transport/
└── examples/
    └── basic/main.go
```

```go
// go.mod
module github.com/your-org/gonice

go 1.22
```

> ⚠️ **`module` path phải khớp exact URL repository.**
> Sai 1 ký tự → `go get` fail `404 Not Found`.

### Bước 2 — Đóng băng API & kiểm tra

```bash
go mod tidy
test -z "$(gofmt -l .)"          # format sạch
goimports -l .
go vet ./...
golangci-lint run
go test ./... -race -cover
govulncheck ./...
go build ./...
```

> 💡 Dùng `gorelease` để phát hiện **breaking change** so với version trước:
> ```bash
> go install golang.org/x/exp/cmd/gorelease@latest
> gorelease   # gợi ý cần MINOR hay MAJOR bump
> ```

### Bước 3 — Gắn tag & push

```bash
git checkout main
git pull

git tag -a v1.0.0 -m "release v1.0.0: initial stable API"
git push origin v1.0.0

# hoặc annotated tag ngắn
git tag v1.1.0 && git push origin v1.1.0
```

### Bước 4 — Kích hoạt Go Module Proxy

```bash
GOPROXY=proxy.golang.org go list -m github.com/your-org/gonice@v1.0.0
go list -m -versions github.com/your-org/gonice      # → v1.0.0 v1.1.0
go install github.com/your-org/gonice/cmd/gonice@latest
```

> 📌 `proxy.golang.org` cache **immutable** — tag đã publish thì
> **không đổi nội dung được**. Nếu tag sai → dùng `retract`, **không** xoá/đổi tag.

### Bước 5 — Phát hành version mới

```bash
# thêm feature → MINOR
git tag v1.1.0 && git push origin v1.1.0

# fix bug → PATCH
git tag v1.1.1 && git push origin v1.1.1
```

**Khi cần breaking change → MAJOR:**

```bash
# 1. Sửa go.mod
#    module github.com/your-org/gonice
#  ▶ module github.com/your-org/gonice/v2
#
# 2. Sửa MỌI import trong repo
#    github.com/your-org/gonice/client
#  ▶ github.com/your-org/gonice/v2/client
#
# 3. Chạy kiểm tra lại
go mod tidy && go test ./... -race

# 4. Tag
git tag v2.0.0 && git push origin v2.0.0
```

### Xử lý tag lỗi — `retract`

```go
// go.mod
retract (
    v1.0.1 // chứa lỗ hổng bảo mật
    v1.0.2 // hồi quy API
)
```

```bash
git commit -am "retract broken releases"
git tag v1.0.3 && git push origin v1.0.3   # release mới để proxy đọc retract
```

> ⚠️ `retract` chỉ là **cảnh báo cho tooling** (`go get` sẽ né version đó).
> Người dùng đã `require` vẫn build được → với CVE nghiêm trọng cần ra **GHSA advisory**.

### Sub-module (version độc lập)

```text
github.com/you/tool/            ← module chính
github.com/you/tool/plugin/     ← module CON (có go.mod riêng, version riêng)
```

> ⚠️ Chỉ tách sub-module khi **thực sự cần** version độc lập —
> tăng độ phức tạp đáng kể (mỗi sub-module 1 chu kỳ release riêng).

---

## 4. Các lệnh quan trọng (CLI)

| Lệnh | Mục đích |
|---|---|
| `go mod init github.com/you/lib` | Tạo module với path chuẩn |
| `go list -m <path>@<version>` | Kiểm tra module/version đã publish chưa |
| `go list -m -versions <path>` | Liệt kê mọi version đã tag |
| `go list -m -json <path>` | Chi tiết module |
| `GOPROXY=proxy.golang.org go list -m <path>@latest` | Kích hoạt / warm proxy |
| `git tag v1.0.0 && git push origin v1.0.0` | Publish release |
| `git tag -d v1.0.0` | Xoá tag **local** (không huỷ được trên proxy) |
| `go install github.com/you/cmd/tool@latest` | Cài CLI tool từ module |
| `gorelease` | So sánh API với version trước → gợi ý major bump |

---

## 5. Lỗi thường gặp & Best Practices

### ⚠️ Các bẫy phổ biến

| # | Bẫy | Hệ quả | Cách tránh |
|---|---|---|---|
| 1 | Tag `1.0.0` (thiếu `v`) | Go không nhận là version | Luôn `vX.Y.Z` |
| 2 | Đổi nội dung tag đã push | `checksum mismatch` cho người dùng | Tag = immutable; dùng tag mới |
| 3 | `module` path ≠ URL repo | `go get` fail 404 | Khớp exactly với VCS path |
| 4 | Bump major quên đổi import path | Module không tồn tại ở path mới | Đổi `go.mod` + **mọi** import |
| 5 | Publish còn `replace` local | Người dùng build fail | Xoá `replace` trước tag |
| 6 | Đẩy thẳng lên `main`, tag nóng vội | Release lỗi lan nhanh | PR + CI xanh + tag từ commit đã review |
| 7 | Không có `CHANGELOG` | Người dùng không biết break cái gì | Ghi mọi release |
| 8 | Export code `internal/` ra public | Mất quyền ẩn, khó refactor về sau | Giữ `internal/` cho unstable code |
| 9 | Xoá tag để "sửa lỗi" | Cache proxy vẫn giữ bản cũ → lệch checksum | **Dùng version mới**, không sửa tag |
| 10 | Bỏ qua LICENSE | Không ai dùng được / vấn đề pháp lý | Có LICENSE ngay từ commit đầu |

### 💡 Best practices

* **Checklist trước tag:**
  ```bash
  go mod tidy -diff      # go.mod sạch
  go vet ./... && golangci-lint run
  go test ./... -race
  govulncheck ./...
  grep -R "replace \.\./" go.mod    # không còn replace local
  ```
* Tag từ **commit đã merge vào `main`**, không từ nhánh local.
* Library quan trọng: chạy `gorelease` trước khi bump major.
* **`v0.x` = API chưa cam kết** — dùng cho giai đoạn dev; khi ổn định mới `v1.0.0`.
* Ưu tiên **Apache-2.0 / MIT** — permisive, dùng được rộng rãi.
* Ghi `CHANGELOG.md` theo [Keep a Changelog](https://keepachangelog.com/).

---

**Bài liên quan**

* [Module and Dependencies](<Module and Dependencies.md>) — `go.mod`, versioning
* [Using Third Party Packages](<Using Third Party Packages.md>) — phía người tiêu dùng dependency
