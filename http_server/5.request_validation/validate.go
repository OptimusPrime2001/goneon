package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
)

func DecodeJson(w http.ResponseWriter, r *http.Request, dst any) bool {
	// 1. Content-Type phải là JSON
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type",
			"Content-Type must be application/json")
		return false
	}

	// 2. Giới hạn kích thước body và từ chối field lạ
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// 3. Decode, rồi phân loại lỗi
	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var maxErr *http.MaxBytesError

		switch {
		case errors.As(err, &maxErr):
			writeError(w, http.StatusRequestEntityTooLarge, "body_too_large", "request body is too large")
		case errors.Is(err, io.EOF):
			writeError(w, http.StatusBadRequest, "empty_body", "request body must not be empty")
		case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
			writeError(w, http.StatusBadRequest, "invalid_json", "request body contains malformed JSON")
		case errors.As(err, &typeErr):
			writeError(w, http.StatusBadRequest, "invalid_type",
				fmt.Sprintf("field %q has the wrong type", typeErr.Field))
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			writeError(w, http.StatusBadRequest, "unknown_field", err.Error())
		default:
			writeError(w, http.StatusBadRequest, "invalid_json", "request body is invalid")
		}
		return false
	}

	// 4. Body chỉ được chứa đúng một giá trị JSON, không được có dữ liệu thừa phía sau
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must contain a single JSON object")
		return false
	}
	return true
}
