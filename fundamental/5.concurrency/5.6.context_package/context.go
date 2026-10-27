package main

// ==========================================
// 1. PACKAGE CONTEXT (context.*)
// ==========================================

// --- Types & Interfaces ---
// context.Context             : Interface chính đại diện cho Context
// context.CancelFunc          : Hàm type func() dùng để chủ động hủy context
// context.CancelCauseFunc     : Hàm type func(cause error) dùng để hủy context kèm lý do

// --- Root Contexts ---
// context.Background()        : Tạo context gốc rỗng (dùng cho main, request top-level)
// context.TODO()              : Tạo context tạm thời khi chưa xác định rõ context nào

// --- Context Constructors (Cancellation & Timeouts) ---
// context.WithCancel()        : Tạo context con có thể hủy thủ công qua CancelFunc
// context.WithCancelCause()   : Tạo context con có thể hủy thủ công kèm lý do (error)
// context.WithTimeout()       : Tạo context tự hủy sau một khoảng thời gian duration
// context.WithTimeoutCause()  : Tạo context tự hủy sau duration kèm lý do chỉ định
// context.WithDeadline()      : Tạo context tự hủy tại mốc thời gian time.Time cụ thể
// context.WithDeadlineCause() : Tạo context tự hủy tại mốc time.Time kèm lý do chỉ định

// --- Utilities & Values ---
// context.WithValue()         : Lưu cặp Key-Value vào context để truyền dữ liệu
// context.WithoutCancel()     : Giữ lại values của context cha nhưng ngắt kết nối hủy/timeout
// context.AfterFunc()         : Đăng ký callback func() tự động chạy ngay khi context bị hủy
// context.Cause()             : Lấy ra error nguyên nhân khiến context bị hủy (từ CancelCauseFunc)

// --- Standard Errors / Variables ---
// context.Canceled            : Lỗi tiêu chuẩn khi context bị hủy thủ công ("context canceled")
// context.DeadlineExceeded    : Lỗi tiêu chuẩn khi context hết hạn thời gian ("context deadline exceeded")

// ==========================================
// 2. CONTEXT INSTANCE (ctx.*)
// ==========================================

// ctx.Done()     : Trả về channel (<-chan struct{}). Channel này sẽ đóng khi context bị hủy/timeout
// ctx.Err()      : Trả về error giải thích lý do ctx kết thúc (nil nếu đang chạy, Canceled, hoặc DeadlineExceeded)
// ctx.Deadline() : Trả về (deadline time.Time, ok bool) cho biết mốc thời gian hết hạn của ctx
// ctx.Value()    : Trả về value kiểu any tương ứng với key truyền vào (ctx.Value(key))
