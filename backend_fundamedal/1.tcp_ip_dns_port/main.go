package main

import (
	"fmt"
	"io"
	"net"
)

func DNS_DISCOVER() {
	ips, _ := net.LookupHost("google.com")
	fmt.Println("A/AAAA", ips)
	cname, _ := net.LookupCNAME("www.github.com")
	fmt.Println("CNAME:", cname)

	mx, _ := net.LookupMX("gmail.com")
	for _, r := range mx {
		fmt.Println("MX:", r.Host, r.Pref)
	}
}

func TCP_DISCOUVER() {
	ln, err := net.Listen("tcp", ":9000") // Lắng nghe port 9000 trên mọi interface
	if err != nil {
		panic(err) // Thử chạy 2 bản cùng lúc để thấy "Address already in use"
	}
	defer ln.Close()

	fmt.Println("Listening on : 9000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			fmt.Println("Client: ", conn.RemoteAddr()) // IP:port ephemeral của client
			io.Copy(conn, conn)                        // đọc từ conn, ghi ngược lại conn
		}(conn)
	}
}
func main() {
	TCP_DISCOUVER()
}
