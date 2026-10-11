# Go Commands — Hướng dẫn chi tiết

Tài liệu tham khảo nhanh về **Go toolchain** (theo `Roadmap.md` §19. Tooling).
Mọi ví dụ đều chạy được trong repo `goneon` này.

> Máy hiện tại: `go version go1.25.3 windows/amd64`
> `GOROOT = C:\Program Files\Go` · `GOPATH = C:\Users\trung\go` · `GOPROXY = https://proxy.golang.org,direct`

---

## 0. Bắt đầu với `go`

### Cú pháp tổng quát

```bash
go <command> [arguments]
```

### Bảng lệnh gốc

| Lệnh | Mục đích |
|-------|----------|
| `run` | Biên dịch và chạy ngay program (không giữ binary) |
| `build` | Biên dịch package + dependency, ghi binary ra thư mục |
| `install` | Biên dịch và cài binary vào `GOBIN`/`GOPATH/bin` |
| `test` | Chạy unit test, benchmark, fuzz |
| `fmt` | Tự động format lại source (gói của `gofmt`) |
| `vet` | Kiểm tra lỗi tiềm ẩn / code đáng nghi |
| `doc` | Xem tài liệu package, hàm, kiểu |
| `list` | Liệt kê package / module |
| `env` | In biến môi trường của Go |
| `mod` | Quản lý module (`init`, `tidy`, `get`…) |
| `work` | Quản lý workspace nhiều module |
| `get` | Thêm dependency và cài (module mode) |
| `clean` | Dọn file object & cache |
| `fix` | Tự động cập nhật code sang API mới |
| `generate` | Chạy lệnh sinh code (`//go:generate`) |
| `version` | In phiên bản Go |
| `tool` | Chạy tool kèm theo (`go tool pprof`, `cover`…) |
| `vet`, `bug`, `telemetry` | Kiểm tra / báo bug / quản lý telemetry |

### Trợ giúp

```bash
go help              # danh sách lệnh
go help <command>    # ví dụ: go help build
go help testflag     # toàn bộ cờ của go test
go help packages     # cách viết pattern package (./..., ./cmd/...)
```

---

## 1. `go env` — kiểm tra môi trường

```bash
go env               # in tất cả
go env GOPATH GOROOT GOFLAGS GOPROXY
go env -json GOPATH  # in dạng JSON
go env GOPATH        # → C:\Users\trung\go
```

| Biến | Ý nghĩa | Giá trị mẫu |
|-------|----------|-------------|
| `GOROOT` | Nơi cài đặt Go | `C:\Program Files\Go` |
| `GOPATH` | Workspace của developer (code, cache, bin) | `C:\Users\trung\go` |
| `GOBIN` | Nơi `go install` ghi binary | mặc định = `$GOPATH/bin` |
| `GOMODCACHE` | Cache module đã tải | `C:\Users\trung\go\pkg\mod` |
| `GOPROXY` | Proxy tải dependency | `https://proxy.golang.org,direct` |
| `GOFLAGS` | Cờ tự động thêm cho mọi lệnh | *(trống)* |
| `GONOSUMCHECK` / `GONOSUMDB` | Bỏ qua checksum (nội bộ, private) | — |
| `GOCACHE` | Cache build | — |

**Mẹo offline / mạng chậm:**

```bash
go env -w GOPROXY=https://goproxy.io,direct   # đổi proxy (lưu vào user config)
go env -w GOPROXY=off                         # hoàn toàn offline (chỉ dùng cache/vendor)
```

---

## 2. Chạy & biên dịch

### 2.1 `go run`

Biên dịch vào RAM/cache rồi chạy ngay — **không** để lại file binary.

```bash
go run main.go        # chạy 1 file
go run .              # chạy package hiện tại (cùng thư mục có go.mod)
go run ./cmd/server   # chạy package con
```

Trong repo này:

```bash
go run main.go
```

> ⚠️ `go run` luôn compile lại (dùng cache build, nhưng vẫn tốn vài trăm ms). Muốn chạy nhanh nhiều lần → build một lần rồi chạy binary.

### 2.2 `go build`

```bash
go build                          # biên dịch package hiện tại
go build ./...                    # biên dịch toàn bộ package trong module
go build -o bin/goneon main.go    # chỉ định tên file output
go build -o dist/ ./cmd/...       # output vào thư mục
```

Trong repo này (giống target `make build`):

```bash
mkdir -p bin
go build -o bin/goneon main.go
./bin/goneon.exe        # Windows
```

**Cờ build quan trọng:**

| Cờ | Ý nghĩa |
|-----|----------|
| `-o <file>` | Tên file kết quả (thư mục hoặc file) |
| `-v` | In tên package đang compile (debug) |
| `-a` | Buộc recompile package đã up-to-date |
| `-race` | Bật data race detection (chỉ `linux/amd64`, `darwin/*`, `windows/amd64`…) |
| `-trimpath` | Bỏ đường dẫn tuyệt đối khỏi binary (reproducible build) |
| `-ldflags="-s -w"` | Bỏ symbol table + DWARF → binary nhỏ hơn ~20–30% |
| `-ldflags="-X 'pkg.var=val'"` | Ghi giá trị biến `string` lúc build (dùng cho version) |
| `-gcflags="-m"` | In escape analysis (xem biến nào bị đẩy ra heap) |
| `-C dir` | `cd` vào thư mục trước khi chạy (phải là cờ đầu tiên) |
| `-n` | Chỉ in lệnh sẽ chạy, không thực hiện |
| `-p n` | Số job build song song (mặc định = số CPU) |

**Build production tối giản:**

```bash
go build -trimpath -ldflags="-s -w" -o bin/goneon main.go
```

**Ghi version vào binary:**

```bash
go build -ldflags="-X main.version=1.2.3" -o bin/goneon main.go
```

### 2.3 `go install`

Giống `go build` nhưng **đưa binary vào `GOBIN`** (`$GOPATH/bin`) để gọi từ mọi nơi.

```bash
go install .                    # cài package hiện tại
go install ./cmd/...            # cài nhiều package
go install golang.org/x/tools/cmd/goimports@latest   # cài tool phiên bản mới nhất
go install golangci-lint/cmd/golangci-lint@v1.61.0   # cài tool phiên bản cố định
```

### 2.4 So sánh `run` / `build` / `install`

| | `go run` | `go build` | `go install` |
|---|---|---|---|
| Chạy ngay | ✅ | ❌ | ❌ |
| Để lại binary | ❌ | ✅ (thư mục hiện tại / `-o`) | ✅ (`GOBIN`) |
| Gọi được từ PATH | ❌ | ❌ (phải `./name`) | ✅ |
| Dùng cho | dev, debug nhanh | release, CI | cài tool |

---

## 3. `go test` — kiểm thử

```bash
go test                    # test các package có thay đổi
go test ./...              # test toàn bộ module
go test -v ./...           # verbose: in từng case
go test -v -race ./...     # như `make test` trong repo này
go test -run TestName ./internal/config/   # chạy 1 test cụ thể
go test -run 'TestConfig/Production' ./internal/config/   # chỉ chạy 1 subtest
go test -count=1 ./...     # bỏ qua cache test
go test -timeout 30s ./... # giới hạn thời gian (mặc định 10m)
```

### Cờ thường dùng

| Cờ | Ý nghĩa |
|----|----------|
| `-v` | In log của mọi test (không chỉ test fail) |
| `-run <regexp>` | Chỉ chạy test match regex (hỗ trợ `Parent/Child`) |
| `-count <n>` | Chạy lại n lần — dùng để **bypass test cache** |
| `-race` | Phát hiện data race (mặc định bật `-covermode=atomic` khi kèm cover) |
| `-timeout <dur>` | Giới hạn thời gian (vd `10m`, `30s`) |
| `-failfast` | Dừng ngay ở test đầu tiên fail |
| `-parallel n` | Số subtest chạy song song (cần `t.Parallel()`) |
| `-short` | Bỏ test dài (dùng `if testing.Short() { t.Skip() }`) |

### Coverage

```bash
go test -cover ./...                         # in % mỗi package
go test -coverprofile=coverage.out ./...     # ghi file profile
go tool cover -func=coverage.out             # % theo từng hàm
go tool cover -html=coverage.out -o cov.html # mở báo cáo HTML
```

Output mẫu:

```text
ok      goneon/internal/config   0.354s  coverage: 87.5% of statements
```

| `-covermode` | Ý nghĩa |
|---------------|----------|
| `set` | Câu lệnh có chạy không (bool) — mặc định |
| `count` | Câu lệnh chạy bao nhiêu lần |
| `atomic` | như `count` nhưng an toàn với goroutine (dùng kèm `-race`) |

### Benchmark

```bash
go test -bench . ./...              # chạy mọi Benchmark
go test -bench=Append -benchmem ./internal/...   # 1 benchmark + số liệu memory
go test -bench=. -benchtime=2s ./...             # kéo dài 2s mỗi benchmark
go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out ./...
go tool pprof cpu.out               # phân tích profile
```

Output mẫu:

```text
BenchmarkAppend-8     1000000    1023 ns/op    416 B/op    4 allocs/op
```

* `B/op` — byte cấp phát mỗi lần gọi, `allocs/op` — số lần `malloc` mỗi lần gọi.
* Mục tiêu: giảm `allocs/op`, không chỉ giảm `ns/op`.

### Fuzz (Go 1.18+)

```bash
go test -fuzz=FuzzParse -fuzztime=30s ./internal/...
go test -fuzz=FuzzX -fuzztime=1x ./...   # chạy đúng 1 seed
```

---

## 4. Format & kiểm tra chất lượng

### 4.1 `go fmt` / `gofmt`

```bash
go fmt ./...          # format toàn bộ package trong module (tương đương gofmt -l -w)
gofmt -l .            # CHỈ LIỆT KÊ file chưa format (dùng cho CI)
gofmt -d main.go      # xem diff trước khi sửa
gofmt -w main.go      # ghi đè file đã format
gofmt -s -w .         # format + simplify code
gofmt -r 'a[b:len(a)] -> a[b:]' -w .   # viết lại theo rule
```

> `go fmt` = `gofmt -l -w` trên các package được chỉ định.

### 4.2 `go vet`

```bash
go vet ./...
go vet -atomic -bool -stringintconv ./...
```

Phát hiện các kiểu lỗi như:

* printf sai format (`%d` với `string`)
* struct field bị copy sai (chứa mutex/pointer)
* unreachable code, cặp `if/else` giống hệt nhau
* method sai signature khi implement interface

> Chạy `go vet` trước `go test` — nó nhanh hơn và bắt lỗi compile-level.

### 4.3 `golangci-lint` (aggregate linter)

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

golangci-lint run                # lint toàn bộ project
golangci-lint run ./internal/... # 1 thư mục
golangci-lint run --fix           # tự sửa lỗi auto-fixable
golangci-lint linters             # liệt kê linter đang bật
golangci-lint fmt                 # format bằng linter formatter
```

`.golangci.yml` cơ bản:

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck      # bắt lỗi return chưa check
    - govet         # go vet
    - staticcheck   # SAxxxx
    - unused
    - ineffassign
    - gosimple
    - misspell

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

---

## 5. Module & dependency

### 5.1 Khởi tạo & dọn dẹp

```bash
go mod init goneon        # tạo go.mod (module name = import path)
go mod tidy               # thêm/bớt dependency khớp với source → đồng bộ go.sum
go mod download           # tải toàn bộ module về cache
go mod verify             # kiểm tra checksum trong go.sum
go mod graph              # in đồ thị dependency
go mod why ./internal/... # vì sao package này lại có trong build?
go mod edit -require=...  # sửa go.mod bằng lệnh
go mod vendor             # copy dependency vào thư mục vendor/
```

`go.mod` trong repo này:

```go
module goneon

go 1.24.5
```

> ⚠️ Luôn chạy `go mod tidy` sau khi thêm/xoá import. Nếu `go.sum` thiếu → lỗi `missing go.sum entry`.

### 5.2 Thêm / nâng cấp dependency — `go get`

```bash
go get github.com/sirupsen/logrus@v1.9.3   # thêm đúng phiên bản
go get github.com/sirupsen/logrus@latest   # bản mới nhất (stable)
go get github.com/sirupsen/logrus@master   # từ nhánh master (không khuyến nghị)
go get -u ./...                            # upgrade mọi dependency dùng trong module
go get -u=false ./...                      # giữ nguyên version
go get -t ./...                            # gồm cả dependency của test
```

**Semantic Import Versioning:**

| Version | Import path |
|---------|-------------|
| `v0.x`, `v1.x` | `github.com/user/pkg` |
| `v2.x` | `github.com/user/pkg/v2` |
| `v3.x` | `github.com/user/pkg/v3` |

Quên `/v2` → lỗi `invalid version: module contains a go.mod file, so module path must match major version (".../v2")`.

### 5.3 Câu lệnh GOFLAGS / thao tác khác

```bash
GOFLAGS=-mod=vendor go build ./...    # build từ vendor/ (offline CI)
go list -m all                        # danh sách module + version
go list ./...                         # danh sách package
go list -m -json all                  # chi tiết dạng JSON
go mod edit -json                     # in go.mod dạng JSON
```

### 5.4 Workspace — `go work` (nhiều module)

```bash
go work init                      # tạo go.work ở thư mục hiện tại
go work use ./api ./web           # thêm module con
go work sync                      # đồng bộ version giữa các module
go work edit -dropuse ./old       # gỡ module
go build ./...                    # build với toàn bộ workspace
```

> `go.work` nên nằm trong `.gitignore` (mỗi dev cấu hình đường dẫn riêng).

---

## 6. Tra cứu & phân tích

### 6.1 `go doc`

```bash
go doc                    # package hiện tại
go doc fmt                # tổng quan package fmt
go doc fmt.Println        # document 1 hàm
go doc -all fmt           # TOÀN BỘ (kể cả private)
go doc -src fmt.Println   # kèm source code
go doc -u strings         # gồm cả symbol không export ở package liên quan
```

### 6.2 `go list`

```bash
go list ./...                          # mọi package trong module
go list -m all                         # mọi module + version
go list -f '{{.ImportPath}} {{.GoFiles}}' ./...   # template tùy biến
go list -deps ./...                    # toàn bộ dependency
```

### 6.3 Escape Analysis

```bash
go build -gcflags="-m" ./...           # biến nào escape ra heap?
go build -gcflags="-m -m" ./...        # chi tiết hơn (2 lần -m)
```

Output mẫu:

```text
./main.go:10:2: moved to heap: buf
./main.go:12:6: leaked param: data
```

### 6.4 `go tool pprof`

```bash
go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out ./...
go tool pprof cpu.out
(pprof) top          # 10 hàm tốn CPU nhất
(pprof) web          # render đồ thị (cần graphviz)
(pprof) list main.go # xem theo dòng code
```

### 6.5 Các lệnh khác

```bash
go clean -cache      # dọn cache build (khi nghi cache hỏng)
go clean -testcache  # chỉ dọn cache test
go fix ./...         # nâng code sang API mới sau khi upgrade Go
go tool              # liệt kê tool đi kèm
```

---

## 7. Cheat sheet

| Việc cần làm | Lệnh |
|--------------|------|
| Chạy nhanh khi dev | `go run .` |
| Build binary | `go build -o bin/goneon main.go` |
| Build release nhỏ gọn | `go build -trimpath -ldflags="-s -w" -o bin/goneon main.go` |
| Chạy test cả project | `go test -v -race ./...` |
| Chạy 1 test | `go test -run TestName -v ./path/` |
| Test + coverage | `go test -coverprofile=c.out ./... && go tool cover -html=c.out` |
| Benchmark | `go test -bench=. -benchmem ./...` |
| Format code | `go fmt ./...` |
| Kiểm tra lỗi | `go vet ./...` |
| Lint | `golangci-lint run --fix` |
| Thêm dependency | `go get pkg@vX.Y.Z` |
| Dọn dependency | `go mod tidy` |
| Xem tài liệu | `go doc fmt.Println` |
| Tra import path | `go list ./...` |

---

## 8. Lỗi thường gặp

| Lỗi | Nguyên nhân | Cách sửa |
|-----|-------------|----------|
| `go: command not found` | Go chưa có trong `PATH` | Cài lại Go, hoặc thêm `%GOROOT%\bin` vào `PATH` |
| `no required module provides package X` | Thiếu dependency / chưa tidy | `go mod tidy` |
| `missing go.sum entry for module` | `go.sum` thiếu checksum | `go mod tidy` (hoặc `go mod download`) |
| `updates to go.mod needed` | go.mod không khớp source | `go mod tidy` |
| `redeclared in this block` | Khai báo trùng | Kiểm tra scope / trùng tên package |
| `undefined: X` (mới `go get`) | IDE chưa index lại | Reload window / `gopls` restart |
| Test pass nhưng **không chạy lại** | **test cache** | `go test -count=1 ./...` |
| `fatal error: race: version mismatch` | binary build thiếu `-race` | Build lại **toàn bộ** bằng `-race` |
| `dial tcp ... i/o timeout` khi `go get` | mạng tới proxy chậm | `go env -w GOPROXY=https://goproxy.io,direct` |
| `build constraints exclude all Go files` | file sai suffix (ví dụ `_windows.go` trên Linux) | Kiểm tra build tags / GOOS |
| `invalid version: ... must match major version` | Thiếu `/v2`, `/v3` trong import path | Sửa import path theo Semantic Versioning |

---

## 9. Áp dụng vào repo `goneon`

Toàn bộ flow phát triển:

```bash
# 1. lần đầu / thêm dependency
go mod tidy

# 2. chạy dev
go run main.go                 # hoặc: make run

# 3. kiểm tra chất lượng
gofmt -l .                     # không được có output
go vet ./...
go test -v -race ./...

# 4. build & chạy binary
go build -trimpath -ldflags="-s -w" -o bin/goneon main.go
./bin/goneon.exe

# 5. chạy thử API
#    GET http://localhost:8080/health
#    GET http://localhost:8080/greet?name=Gopher
```

Các biến môi trường hỗ trợ: `PORT` (mặc định `8080`), `ENV` (`development`/`production`/`test`), `LOG_LEVEL` (`debug`/`info`/`warn`/`error`).

```bash
PORT=9000 LOG_LEVEL=debug go run main.go     # PowerShell: $env:PORT=9000; $env:LOG_LEVEL=debug; go run main.go
```

---

**Tham khảo thêm:** `go help <command>` · `go help testflag` · [pkg.go.dev](https://pkg.go.dev)
