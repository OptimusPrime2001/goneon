package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func RequestBasic() {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get("https://example.com")

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println("Status:", response.Status)                     // 200 OK
	fmt.Println("Proto:", response.Proto)                       // HTTP/2.0
	fmt.Println("TLS version:", response.TLS.Version == 0x0304) // true nếu TLS 1.3
	fmt.Println("Cert subject:", response.TLS.PeerCertificates[0].Subject)
	fmt.Println("Cert expires:", response.TLS.PeerCertificates[0].NotAfter)

	body, _ := io.ReadAll(response.Body)
	fmt.Println("Body size:", len(body))
}
func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s from %s\n", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)
		fmt.Fprintln(w, "Xin chào từ Go!")
	})
	fmt.Println("listening on :8080")
	if err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil); err != nil {
		panic(err)
	}
}
