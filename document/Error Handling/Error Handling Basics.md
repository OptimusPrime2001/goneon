# Error Handling Basics trong Go

Go xử lý lỗi theo cách **rất khác** Java/Python: không có `try/catch`/`throw`.
Lỗi trong Go là **một giá trị bình thường** được trả về, và việc kiểm tra nó là
**control flow tường minh** do bạn tự viết.

Khớp node **Error Handling** của roadmap.sh (§7).
Code mẫu: `fundamental/4.error_handling/error_basic.go`

---

## 1. Khái niệm & Bản chất

> **Triết lý cốt lõi:** *"Error is a value."* — Lỗi không phải ngoại lệ (exception)
> mà là **giá trị** bạn trả về và xử lý như mọi giá trị khác.

Điều này dẫn đến 3 quy tắc vàng của Go:

1. **Không có `try/catch`.** Hàm nào có thể thất bại thì **tự trả về `error`**.
2. **`error` luôn là tham số trả về *cuối cùng*.** Người gọi **bắt buộc** nhìn thấy nó.
3. **Kiểm tra `err != nil` ngay lập tức** — fail fast, đừng để lẫn sang câu lệnh sau.

```text
  Exception-based (Java/Python)              Lập trình Go
  ─────────────────────────────────        ─────────────────────────────
  try {                                   result, err := DoWork()
      result = DoWork()                   if err != nil {           // kiểm tra NGAY
  } catch e {                                return err              // hoặc xử lý
      handle(e)                           }
  }                                       use(result)
```

> **Vì sao Go chọn cách này?** Compiler Go **không** cho phép bỏ qua giá trị trả về.
> Nếu hàm trả `error` mà bạn không gán/xử lý → `go vet` báo lỗi.
> Nhờ đó **không thể vô tình nuốt lỗi** như các ngôn ngữ exception.

---

## 2. Cú pháp chuẩn: `(T, error)`

Một hàm có thể thất bại sẽ trả về **cặp** `(giá trị hữu ích, error)`:

```go
// Divide trả về (kết quả, error). Khi thất bại: (zero value, lỗi).
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("divide: chia cho 0") // (zero value, error)
    }
    return a / b, nil // (value, nil) = thành công
}
```

### Quy tắc: `nil` nghĩa là thành công

Trong Go, convention là **`error == nil` ⇔ thành công** (đảo ngược so với
`Result`/`Optional` của các ngôn ngữ khác):

```text
err == nil   →  THÀNH CÔNG, dùng được giá trị trả về
err != nil   →  THẤT BẠI,  giá trị còn lại thường là zero value (0, "", nil...)
```

### Gọi và kiểm tra — fail fast

```go
if result, err := Divide(10, 2); err != nil {
    fmt.Println("Lỗi:", err)
} else {
    fmt.Println("10 / 2 =", result)
}

if _, err := Divide(10, 0); err != nil {
    fmt.Println("Lỗi:", err)
}
```

Output thật:

```text
10 / 2 = 5
Lỗi: divide: chia cho 0
```

> Pattern `if result, err := f(); err != nil { ... } else { ... }` gọn và chuẩn.
> Nếu không muốn dùng `else`, có thể **return sớm** (early return) cho dễ đọc.

---

## 3. Đọc error từ hàm chuẩn library

Ngay cả stdlib cũng trả `error` — bạn phải kiểm tra:

```go
f, err := os.Open("config.txt")
if err != nil {
    log.Fatal(err) // hoặc return err
}
defer f.Close()
```

```go
data, err := os.ReadFile("khong_ton_tai.txt")
fmt.Println("err != nil ?", err != nil)
fmt.Println("err =", err)
```

Output thật:

```text
err != nil ? true
err = open khong_ton_tai.txt: no such file or directory
```

> Hàm như `os.Open`, `os.ReadFile`, `http.Get`, `json.Unmarshal`... **luôn** trả `error`.
> Bỏ qua nó (`f, _ := os.Open(...)`) là lỗi programming nghiêm trọng.

---

## 4. 3 sai lầm kinh điển cần tránh

```text
Sai lầm                          Vì sao sai
───────────────────────────────  ─────────────────────────────────────────────
① Bỏ qua lỗi:  f, _ := Open()   Nuốt lỗi im lặng → bug khó tìm, data corruption
② Log rồi nuốt: log.Println(e)  "Xử lý" nhưng KHÔNG return → caller tưởng thành công
                                 return nil
③ So sánh chuỗi:                Vỡ khi lỗi bị wrap/thêm ngữ cảnh ở tầng khác
   strings.Contains(err.Error())
```

### Cách đúng: thêm ngữ cảnh rồi **trả lên** (propagate)

```go
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // %w bọc lỗi gốc + thêm ngữ cảnh "đang làm gì"
        return nil, fmt.Errorf("LoadConfig(%s): %w", path, err)
    }
    // ...
    return cfg, nil
}
```

> **Quy tắc:** tầng nào **hiểu** lỗi thì xử lý; tầng nào **không hiểu** thì
> **bọc ngữ cảnh và trả lên**. Đừng bao giờ `return nil` sau khi chỉ log.

---

## 5. Best Practices / Lưu ý quan trọng

* **Error là tham số trả về cuối cùng** — `(T, error)`, không có ngoại lệ.
* **`err == nil` = thành công** — kiểm tra `if err != nil` ngay sau mỗi lời gọi.
* **Fail fast** — trả về ngay khi gặp lỗi, đừng chạy tiếp với giá trị không hợp lệ.
* **Không bỏ qua lỗi** — đừng dùng `_` để nuốt `error`; chỉ `_` khi bạn **thực sự**
  chứng minh được lỗi đó vô hại (hiếm).
* **Thêm ngữ cảnh khi propagate** — dùng `fmt.Errorf("...: %w", err)` để giữ nguyên
  lỗi gốc (hỗ trợ `errors.Is`/`errors.As` ở bài sau).
* **Không nuốt lỗi bằng log + return nil** — hoặc xử lý triệt để, hoặc trả lên cho caller.
* **Viết message lỗi lower-case**, không dấu chấm cuối, wrap ngữ cảnh ở mỗi tầng.

---

## 6. Summary

* Go **không có exception** — lỗi là **giá trị `error`** được trả về tường minh.
* `error` luôn là **tham số trả về cuối cùng**, theo cặp `(T, error)`.
* **`err == nil` = thành công**; kiểm tra `if err != nil` ngay, fail fast.
* Hàm chuẩn library (`os.Open`...) **luôn** trả `error` — đừng bao giờ bỏ qua (`_`).
* Tránh 3 sai lầm: bỏ qua lỗi, log-roi-nuốt, so sánh chuỗi `err.Error()`.
* Khi propagate: **bọc ngữ cảnh bằng `%w`** và trả lên cho tầng hiểu lỗi xử lý.
