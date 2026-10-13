# Methods trong Go (Golang)

Method là hàm "gắn" vào một **named type** — cú pháp thêm *receiver*: `func (r T) Name()`.

Khớp node **Methods** của roadmap.sh (§5).
Code mẫu: `fundamental/4.method_interface/method.go`

---

## 1. Method Declaration

```go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {   // (r Rectangle) = RECEIVER
    return r.Width * r.Height
}
```

* `(r Rectangle)` là **receiver** — method gắn với type `Rectangle`
* Gọi như hàm thường: `r.Area()` — biến `r` tự động được truyền làm receiver

### Receiver phải là named type

```go
type Celsius float64
func (c Celsius) ToFahrenheit() float64 { ... }  // ✅ — type Celsius "tên thật"

type C = Celsius                                  // alias
func (c C) Foo() {}                               // ❌ compile error
// cannot define new methods on alias type
```

---

## 2. Value Receiver — làm việc trên BẢN COPY

```go
type Account struct{ Balance int }

// Deposit (value receiver): a.Balance += amount chỉ áp dụng trên BẢN COPY.
func (a Account) Deposit(amount int) {
    a.Balance += amount
}

acc := Account{Balance: 100}
acc.Deposit(50)

fmt.Println(acc.Balance)     // 100 — KHÔNG đổi!
```

Output thật:

```text
Deposit value receiver — Balance: 100
```

* Go luôn pass by value — receiver cũng được **copy** mỗi lần gọi
* Thích hợp cho method **chỉ đọc**: getter, `Area()`, `String()`

---

## 3. Pointer Receiver — sửa trên BẢN GỐC

```go
func (a *Account) DepositReal(amount int) {
    a.Balance += amount      // sửa trực tiếp trên bản gốc
}

acc := Account{Balance: 100}
acc.DepositReal(50)

fmt.Println(acc.Balance)     // 150
```

Output thật:

```text
Deposit pointer receiver — Balance: 150
```

### Auto-address & auto-deref

Compiler tự "đi qua" con trỏ — hiếm khi cần ghi tường minh:

```go
acc.DepositReal(50)          // acc là value  → tự thành (&acc).DepositReal(50)

p := &acc
fmt.Println(p.Balance)       // p là pointer → tự thành (*p).Balance
```

---

## 4. Method Sets ⭐ (quan trọng khi dùng interface)

**Method set** = tập hợp các method mà type **thực sự "có"** để satisfy interface:

| Receiver khai báo trên | Method set của `T` | Method set của `*T` |
|---|---|---|
| Value `(T)` | ✅ có | ✅ có |
| Pointer `(*T)` | ❌ **không** | ✅ có |

Ví dụ thấy rõ:

```go
type Scaler interface{ Scale(f float64) }

// Rectangle.Scale nhận POINTER receiver:
var sc Scaler = Rectangle{1, 2}      // ❌ compile error
// cannot use Rectangle{…} (value of type Rectangle) as Scaler value:
// Rectangle does not implement Scaler (method Scale requires pointer receiver)

sc = &Rectangle{1, 2}                // ✅ *Rectangle có method set đầy đủ
```

**Kết luận nhanh:**

* Type `T` chỉ satisfy interface nếu interface chỉ cần method **value receiver**
* Type `*T` satisfy interface cần **cả** value lẫn pointer receiver methods

---

## 5. Method vs Function

| | Method | Function |
|---|---|---|
| Gắn với type? | ✅ — qua receiver | ❌ — tự do |
| Cách gọi | `r.Method()` | `f(args)` |
| Tham gia interface? | ✅ — qua method set | ❌ |
| Tên giống nhau giữa các type? | ✅ — mỗi type1 tên riêng | ❌ — trùng tên = lỗi |

### Method expression — biến method thành hàm

```go
areaFn := Rectangle.Area      // Rectangle.Area là hàm: receiver thành tham số đầu
fmt.Println(areaFn(r))        // 200
```

Output thật:

```text
Method expression: 200
```

Dùng khi cần truyền method như callback, hoặc muốn gọi tường minh kiểu nào đứng đầu.

---

## 6. Khi nào dùng receiver nào?

| Tình huống | Chọn |
|---|---|
| Method chỉ đọc, không mutate | Value |
| Cần mutate gốc (setter, `Deposit`…) | Pointer |
| Struct lớn (tránh copy) | Pointer |
| Type cần satisfy interface có method mutate | Pointer |

**Quy tắc thực tế của cộng đồng Go:**

* Struct **nhỏ**, receiver **chỉ đọc** → value
* Đã có **bất kỳ** pointer receiver method nào trên type → dùng pointer cho **tất cả** method còn lại (tránh confusion)
* Struct lớn → pointer

---

## Summary

* Method = hàm + receiver: `func (r T) Name()` — receiver phải là **named type**, không phải alias
* **Value receiver**: copy — không mutate gốc; **pointer receiver**: sửa trên bản gốc
* Compiler tự auto-address (`r.M()` ≡ `(&r).M()`) và auto-deref (`p.field` ≡ `(*p).field`)
* **Method sets**: `T` không có method pointer receiver — `*T` có **đầy đủ** → quyết định type nào satisfy interface
* Method vs Function: method gắn type & tham gia interface; method expression `T.Method(t)` biến method thành hàm thường
* Thực tế: đã có pointer method → dùng pointer hết cho type đó
