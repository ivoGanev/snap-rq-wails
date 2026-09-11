package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// EchoResponse is what the echo server sends back.
type EchoResponse struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// enableCORS adds permissive CORS headers so the mock server can be hit from
// any origin during manual client testing.
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}

	resp := EchoResponse{
		Method:  r.Method,
		URL:     r.URL.String(),
		Headers: headers,
		Body:    string(body),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ---------- Content-type mocks ----------

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := map[string]any{
		"status":  "ok",
		"message": "This is a JSON response",
		"data": map[string]any{
			"id":   42,
			"name": "mock-item",
			"tags": []string{"alpha", "beta", "gamma"},
		},
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func csvHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("id,name,status\n1,alpha,active\n2,beta,pending\n3,gamma,archived\n"))
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Mock HTML</title></head>
<body>
  <h1>Hello from the mock server</h1>
  <p>This is a simple HTML snippet for testing the response viewer.</p>
  <ul>
    <li>Item one</li>
    <li>Item two</li>
  </ul>
</body>
</html>`))
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("This is a plain text response from the mock server.\n"))
}

func xmlHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<root>
  <status>ok</status>
  <item id="1">
    <name>alpha</name>
    <active>true</active>
  </item>
  <item id="2">
    <name>beta</name>
    <active>false</active>
  </item>
</root>`))
}

func binaryHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=mock.bin")
	w.WriteHeader(http.StatusOK)
	// 64 bytes of pseudo-binary data.
	data := make([]byte, 64)
	for i := range data {
		data[i] = byte(i)
	}
	_, _ = w.Write(data)
}

// ---------- Status mocks ----------

func statusHandler(w http.ResponseWriter, r *http.Request, status int, body []byte) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(body) > 0 {
		_, _ = w.Write(body)
	}
}

func status200Handler(w http.ResponseWriter, r *http.Request) {
	statusHandler(w, r, http.StatusOK, []byte(`{"status":"ok","message":"200 OK"}`))
}

func status201Handler(w http.ResponseWriter, r *http.Request) {
	statusHandler(w, r, http.StatusCreated, []byte(`{"status":"created","message":"201 Created"}`))
}

func status204Handler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.WriteHeader(http.StatusNoContent)
}

func status400Handler(w http.ResponseWriter, r *http.Request) {
	statusHandler(w, r, http.StatusBadRequest, []byte(`{"error":"bad_request","message":"400 Bad Request"}`))
}

func status401Handler(w http.ResponseWriter, r *http.Request) {
	statusHandler(w, r, http.StatusUnauthorized, []byte(`{"error":"unauthorized","message":"401 Unauthorized"}`))
}

func status404Handler(w http.ResponseWriter, r *http.Request) {
	statusHandler(w, r, http.StatusNotFound, []byte(`{"error":"not_found","message":"404 Not Found"}`))
}

func status500Handler(w http.ResponseWriter, r *http.Request) {
	statusHandler(w, r, http.StatusInternalServerError, []byte(`{"error":"internal_error","message":"500 Internal Server Error"}`))
}

// ---------- Special mocks ----------

func delayHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	ms := 1000
	if raw := r.URL.Query().Get("ms"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed >= 0 {
			ms = parsed
		}
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":      "ok",
		"delayed_ms":  ms,
		"message":     fmt.Sprintf("Response delayed by %dms", ms),
	})
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`{"error":"not_found","message":"Unknown mock endpoint"}`))
}

func main() {
	mux := http.NewServeMux()

	// Original echo endpoint, preserved for compatibility.
	mux.HandleFunc("/echo", echoHandler)

	// Content-type mocks.
	mux.HandleFunc("/mock/json", jsonHandler)
	mux.HandleFunc("/mock/csv", csvHandler)
	mux.HandleFunc("/mock/html", htmlHandler)
	mux.HandleFunc("/mock/text", textHandler)
	mux.HandleFunc("/mock/xml", xmlHandler)
	mux.HandleFunc("/mock/binary", binaryHandler)

	// Status mocks.
	mux.HandleFunc("/mock/status/200", status200Handler)
	mux.HandleFunc("/mock/status/201", status201Handler)
	mux.HandleFunc("/mock/status/204", status204Handler)
	mux.HandleFunc("/mock/status/400", status400Handler)
	mux.HandleFunc("/mock/status/401", status401Handler)
	mux.HandleFunc("/mock/status/404", status404Handler)
	mux.HandleFunc("/mock/status/500", status500Handler)

	// Special mocks.
	mux.HandleFunc("/mock/delay", delayHandler)
	mux.HandleFunc("/mock/empty", emptyHandler)

	// Redirect root to a tiny help page.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFoundHandler(w, r)
			return
		}
		enableCORS(w)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`Mock server endpoints:

Content-type mocks:
  GET /mock/json    -> application/json
  GET /mock/csv     -> text/csv
  GET /mock/html    -> text/html
  GET /mock/text    -> text/plain
  GET /mock/xml     -> application/xml
  GET /mock/binary  -> application/octet-stream

Status mocks:
  GET /mock/status/200 -> 200 OK
  GET /mock/status/201 -> 201 Created
  GET /mock/status/204 -> 204 No Content
  GET /mock/status/400 -> 400 Bad Request
  GET /mock/status/401 -> 401 Unauthorized
  GET /mock/status/404 -> 404 Not Found
  GET /mock/status/500 -> 500 Internal Server Error

Special mocks:
  GET /mock/delay?ms=2000 -> delayed response (default 1000ms)
  GET /mock/empty           -> 200 OK with empty body

Echo:
  ANY /echo -> echoes request back as JSON
`))
	})

	port := "18080"
	fmt.Printf("Mock server listening on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
