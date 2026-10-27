package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func QuyTrinhHandleFile() {
	filePath := "source/test.txt"
	content := "Dòng dữ liệu mới được thêm vào\n"

	// 1. Mở file để vừa ĐỌC vừa GHI THÊM (Append), tự tạo file nếu chưa có
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Println("Lỗi mở/tạo file:", err)
		return // Thoát ngay nếu không mở được file
	}
	defer file.Close() // Đảm bảo đóng file khi kết thúc hàm main

	// Ghi thêm dữ liệu vào file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Lỗi ghi file:", err)
		return
	}
	fmt.Println("Đã ghi thêm dữ liệu thành công!")

	// 2. Tua con trỏ về ĐẦU file (offset 0) để chuẩn bị ĐỌC
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Lỗi tua con trỏ file:", err)
		return
	}

	// Đọc từng dòng
	fmt.Println("--- NỘI DUNG FILE ---")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("Đọc dòng:", scanner.Text())
	}
	// Kiểm tra lỗi đọc sau khi kết thúc vòng lặp Scan
	if err := scanner.Err(); err != nil {
		fmt.Println("Lỗi khi đọc file:", err)
		return
	}

	// 3. Tua con trỏ về ĐẦU file lần nữa để COPY
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Lỗi tua con trỏ file:", err)
		return
	}

	// Chép dữ liệu qua file khác
	dst, err := os.Create("source/copy.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file copy:", err)
		return
	}
	defer dst.Close()

	bytesCopied, err := io.Copy(dst, file)
	if err != nil {
		fmt.Println("Lỗi khi copy file:", err)
		return
	}

	fmt.Printf("Đã copy thành công %d bytes!\n", bytesCopied)
}
func main() {
	QuyTrinhHandleFile()
}
