# 📚 Tài liệu Go — Mục lục

Thứ tự học lấy theo `Roadmap.md` → mục *Recommended Learning Order*.
Folder **không đánh số** (để khỏi phải đổi tên khi thêm chương) → thứ tự học nằm **tại đây**.

---

## Phân nhóm kiểu dữ liệu trong Go

Go spec (§Types) chỉ chia **2 ô lớn**:

| Ô lớn | Thành phần |
|-------|-----------|
| **Predeclared / basic types** | `bool`, numeric (`int`, `float64`…), `string` |
| **Composite types** ⭐ | `array`, `struct`, `pointer`, `function`, `interface`, `slice`, `map`, `channel` |

> 📖 Nguyên văn spec: *"Composite types — array, struct, pointer, function, interface, slice, map, and channel types — may be constructed using type literals."*
> Từ **"aggregate" không có trong Go spec** — đây là thuật ngữ đi dạy học, chia `composite` theo ngữ nghĩa khi gán giá trị:

```text
Composite types  (Go spec — ô lớn, 8 kiểu)
├── Aggregate  → Array, Struct            value semantics: gán = COPY toàn bộ giá trị
├── Reference  → Slice, Map, Pointer,     có con trỏ bên trong: gán = dùng chung
│                Function, Channel                underlying data
└── Interface  → Interface                "self-contained OR contain references"
                                          → kiếp trung gian, học riêng (Roadmap §6)
```

**Tóm tắt bằng 1 dòng:** `composite` trả lời *"cái gì không phải basic type"*,
`aggregate/reference` trả lời *"gán vào biến thì **copy** hay **share**"*.

`document/Composite Types/` để **phẳng** (mỗi kiểu 1 file), cột **Nhóm** trong bảng dưới đây giữ vai trò của ranh giới copy-vs-share:

```text
document/Composite Types/
├── Array.md      ← Aggregate  (gán = COPY toàn bộ)
├── Struct.md     ← Aggregate  (gán = COPY toàn bộ)
├── Slice.md      ← Reference  (gán = SHARE underlying data)
├── Map.md        ← Reference  (gán = SHARE underlying data)
├── Functions.md  ← composite type; closure capture theo reference
├── String.md     ← ghi chú riêng của bạn (Bytes / ASCII / UTF-8)
└── Pointer.md, Channel.md, Interface.md  → chưa viết
```

---

## 1. Basic Review — Ngôn ngữ nền

### Data types

| Bài | Ghi chú |
|-----|---------|
| [Data type](<Basic Review/Data types/Data type.md>) | Zero value, type inference, conversion |
| [Boolean](<Basic Review/Data types/Boolean.md>) | |
| [Numeric types](<Basic Review/Data types/Numeric types.md>) | |
| [Runes](<Basic Review/Data types/Runes.md>) | `rune` = `int32`, Unicode code point |
| [Strings](<Basic Review/Data types/Strings.md>) | `immutable`, UTF-8 |

> ⚠️ **Tại sao `Strings` nằm ở `Basic Review` chứ không trong `Composite Types/`?**
> Go spec xếp `string` là *predeclared type*, **không** phải composite type.
> Nhưng nó **hành vi như array chỉ-đọc**: có `s[i]`, `s[a:b]`, `len(s)`, `range`.
> → Đọc [Slice](<Composite Types/Slice.md>) và [Array](<Composite Types/Array.md>) song song để thấy điểm tương đồng.

### Variable & Constant

| Bài |
|-----|
| [Variable & Constant](<Basic Review/Variable & Constant/Variable & Constant.md>) |
| [Scope & Shadowing](<Basic Review/Variable & Constant/Scope & Shadowing.md>) |

---

## 2. Composite Types — 8 kiểu xây bằng type literal

Folder `Composite Types/` để **phẳng**; cột **Nhóm** phân biệt ngữ nghĩa khi gán giá trị.

### 2.1 Aggregate — *gán = copy toàn bộ giá trị*

| Bài | Code mẫu |
|-----|-----------|
| [Array](<Composite Types/Array.md>) | `fundamental/2.composite_type/array.go` |
| [Struct](<Composite Types/Struct.md>) | `fundamental/2.composite_type/struct.go` |

> **Đặc điểm:** giá trị nằm **trực tiếp** trong biến → `b := a` tạo **bản sao hoàn toàn**.
> Sửa `b` không ảnh hưởng `a`. Array còn có **length là một phần của kiểu** → `[3]int ≠ [4]int`.

### 2.2 Reference — *gán = dùng chung dữ liệu gốc*

| Bài | Code mẫu | Trạng thái |
|------|-----------|-----------|
| [Slice](<Composite Types/Slice.md>) | `fundamental/2.composite_type/slice.go` | ✅ |
| [Map](<Composite Types/Map.md>) | `fundamental/2.composite_type/map.go` | ✅ |
| [Functions](<Composite Types/Functions.md>) | `fundamental/2.composite_type/loop_function.go` | ✅ |
| Pointer *(Roadmap §3)* | — | ⬜ chưa viết |
| Channel *(Roadmap §11)* | — | ⬜ chưa viết |

> **Đặc điểm:** biến chỉ giữ **descriptor/con trỏ** tới dữ liệu gốc → `b := a` chỉ copy "cái đầu trỏ",
> `b` và `a` **dùng chung** underlying data → sửa `b` làm thay đổi `a` (với slice/map).

### 2.3 Kiếp trung gian & ghi chú

| Bài | Trạng thái |
|-----|-----------|
| `String.md` — ghi chú Byte / ASCII / UTF-8 của bạn | ✅ |
| **Interface** *(Roadmap §6)* | ⬜ chưa viết |

> `interface` là composite type nhưng **không thuần aggregate lẫn reference**
> (spec: *"an interface value may be self-contained **or** contain references to underlying data"*).
> `string` thì Go xếp vào *predeclared type* — bài viết ở [Strings](<Basic Review/Data types/Strings.md>) (mục 1).

---

## 📦 Code Organization — Packages & Modules

| Bài | Nội dung |
|-----|----------|
| [Module and Dependencies](<Code Organization/Module and Dependencies.md>) | `go.mod`, `go.sum`, SemVer, MVS, `replace`/`retract`, `go.work` |
| [Packages](<Code Organization/Packages.md>) | Cấu trúc package, naming, `cmd/` `internal/` `pkg/`, visibility, `init()` |
| [Package Import Rules](<Code Organization/Package Import Rules.md>) | Thứ tự import, alias/blank/dot import, chống circular |
| [Using Third Party Packages](<Code Organization/Using Third Party Packages.md>) | Checklist chọn dependency, pin version, `GOPRIVATE`, vendor, audit |
| [Publishing Modules](<Code Organization/Publishing Modules.md>) | Publish qua Git tag, Semantic Import Versioning, `retract`, `gorelease` |

> Thứ tự đọc: **Module and Dependencies → Packages → Package Import Rules → Using Third Party → Publishing Modules**

---

## 🛠 Tooling

| Bài | Nội dung |
|-----|----------|
| [Commands](<Tooling/Commands.md>) | Go toolchain chi tiết: `go build`, `go test`, `go vet`, `go mod`, `golangci-lint`… |

---

## 🛡 Error Handling — Xử lý lỗi

Folder `Error Handling/` để **phẳng**; đọc theo thứ tự từ nền tảng → nâng cao.

| Bài | Nội dung chính |
|-----|----------------|
| [Error Handling Basics](<Error Handling/Error Handling Basics.md>) | Triết lý "error là giá trị", idiomatic `(T, error)`, `err != nil`, fail fast, tránh nuốt lỗi |
| [The error Interface](<Error Handling/The error Interface.md>) | Bản chất interface `error{ Error() string }`, custom error struct, bẫy **typed-nil** |
| [Creating Errors](<Error Handling/Creating Errors.md>) | `errors.New`, `fmt.Errorf`, khác biệt `%w` vs `%v`/`%s`, nhiều `%w` (Go 1.20+) |
| [Error Wrapping](<Error Handling/Error Wrapping.md>) | `errors.Unwrap`, `errors.Is`, `errors.As`, method `Is`/`Unwrap` tùy biến, `errors.Join` |
| [Sentinel Errors](<Error Handling/Sentinel Errors.md>) | Biến error export sẵn (`ErrXxx`), `errors.Is` vs `==`, sentinel stdlib (`io.EOF`, `os.ErrNotExist`) |
| [panic and recover](<Error Handling/panic and recover.md>) | `panic` (unwind + LIFO defer), `recover` trong `defer`, runtime panic, khi nào panic hợp lệ |
| [Stack Traces and Debugging](<Error Handling/Stack Traces and Debugging.md>) | `runtime/debug.Stack`, `GOTRACEBACK`, đọc/hiểu stack trace, mẹo debug |

> Thứ tự đọc: **Basics → error Interface → Creating Errors → Error Wrapping → Sentinel → panic and recover → Stack Traces**
> Code mẫu: `fundamental/4.error_handling/`

---

## 🔄 Concurrency — Đồng thời (Roadmap §11)

Folder `Concurrency/` để **phẳng**; đọc từ nền tảng → nâng cao.

| Bài | Nội dung chính | Trạng thái |
|-----|----------------|-----------|
| [Goroutines](<Concurrency/Goroutines.md>) | Goroutine là gì, `go` keyword, lifecycle, concurrency ≠ parallelism, đồng bộ bằng channel, data race, goroutine leak, pitfall (sleep/loop capture/panic) | ✅ |
| [Channels](<Concurrency/Channels.md>) | Unbuffered/buffered, directional channel, `close`, `range`, comma ok, **Select** (non-block, timeout, nil channel) | ✅ |
| [Select](<Concurrency/Select.md>) | ❌ **Đã gộp vào `Channels.md`** (mục 5) | ❌ gộp |
| [Worker Pool](<Concurrency/Worker Pool.md>) | `sync.WaitGroup`, worker pool cố định, job động, thứ tự close, chống deadlock | ✅ |
| [Synchronization](<Concurrency/Synchronization.md>) | Mutex, RWMutex, Atomic | ⬜ |
| [Context](<Concurrency/Context.md>) | Cancel, timeout, truyền context | ⬜ |
| [Concurrency Patterns](<Concurrency/Concurrency Patterns.md>) | Fan-Out/Fan-In, Pipeline, Producer-Consumer | ⬜ |
| Common Problems *(gộp vào bài trên)* | Race Condition & Goroutine Leak → `Goroutines.md`; Deadlock → `Worker Pool.md` | ✅ |

> Thứ tự đọc: **Goroutines → Channels (+ Select) → Worker Pool (+ WaitGroup) → Synchronization → Context → Concurrency Patterns**
> Code mẫu: `fundamental/5.concurrency/`

---

## Các file còn lại

| File | Nội dung |
|------|----------|
| [overview.md](overview.md) | Bản sao của `Roadmap.md` (bản gốc ở gốc repo) |
| [keyword.md](keyword.md) | Giải nghĩa thuật ngữ nhanh |

---

## Roadmap chưa viết thành doc

`Roadmap.md` có 21 chương; `document/` mới có chương **1, 2, 3** (gộp vào `Basic Review/` và `Composite Types/`) + **7** (`Error Handling/`) + **8** (`Code Organization/`) + **11 (một phần — còn `Synchronization`/`Context`/`Concurrency Patterns`)** + `Tooling/`.

* 4. Control Flow · 5. Methods · 6. Interfaces
* 9. Generics · 10. Memory Model · 12. Runtime Internals
* 13. Standard Library · 14. File Handling · 15. Testing
* 16. HTTP Server · 17. Database · 18. Profiling · 20. Advanced · 21. Production

---

**Quy ước đặt tên**

* Folder cấp 1: `Title Case`, **không đánh số** (thứ tự do file mục lục này quy định).
* `Composite Types/` để **phẳng** — mỗi kiểu 1 file, không tạo thư mục con.
* File bài học: `Title Case.md` (`Array.md`, `Slice.md`, `Map.md`, `Strings.md`).
* Mỗi bài: khái niệm → cú pháp → ví dụ `go` → output `text` → lỗi thường gặp.
* Code mẫu đặt trong `fundamental/2.composite_type/`, đặt tên trùng với chủ đề bài (`map.go`, `slice.go`…).
