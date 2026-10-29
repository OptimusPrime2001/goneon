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

func writeDecodeError(w http.ResponseWriter, err error) {
	var maxErr *http.MaxBytesError
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError

	switch {
	case errors.As(err, &maxErr):
		writeError(
			w,
			http.StatusRequestEntityTooLarge,
			"body_too_large",
			"request body is too large",
		)

	case errors.Is(err, io.EOF):
		writeError(
			w,
			http.StatusBadRequest,
			"empty_body",
			"request body must not be empty",
		)

	case errors.As(err, &syntaxErr),
		errors.Is(err, io.ErrUnexpectedEOF):
		writeError(
			w,
			http.StatusBadRequest,
			"invalid_json",
			"request body contains malformed JSON",
		)

	case errors.As(err, &typeErr):
		writeError(
			w,
			http.StatusBadRequest,
			"invalid_type",
			fmt.Sprintf("field %q has the wrong type", typeErr.Field),
		)

	case strings.HasPrefix(err.Error(), "json: unknown field"):
		writeError(
			w,
			http.StatusBadRequest,
			"unknown_field",
			"request contains an unknown field",
		)

	default:
		writeError(
			w,
			http.StatusBadRequest,
			"invalid_json",
			"request body is invalid",
		)
	}
}
func decodeJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) bool {
	// 1. Validate Content-Type.
	contentType := r.Header.Get("Content-Type")

	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil || mt != "application/json" {
		writeError(
			w,
			http.StatusUnsupportedMediaType,
			"unsupported_media_type",
			"Content-Type must be application/json",
		)
		return false
	}

	// 2. Limit request body to 1 MiB.
	if r.Body == nil {
		writeError(
			w,
			http.StatusBadRequest,
			"empty_body",
			"request body must not be empty",
		)
		return false
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// 3. Decode the first JSON value.
	if err := dec.Decode(dst); err != nil {
		writeDecodeError(w, err)
		return false
	}

	// 4. Reject a second JSON value or trailing invalid data.
	var extra json.RawMessage

	if err := dec.Decode(&extra); err != nil {
		if errors.Is(err, io.EOF) {
			return true
		}

		writeDecodeError(w, err)
		return false
	}

	// A second valid JSON value was found.
	writeError(
		w,
		http.StatusBadRequest,
		"multiple_json_values",
		"request body must contain only one JSON value",
	)
	return false
}
