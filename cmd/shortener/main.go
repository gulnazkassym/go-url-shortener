package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mu       = sync.Mutex{}
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", postHandler)
	mux.HandleFunc("/hevfyegruf", getHandler)
	fmt.Println("Server started at http://localhost:8080")

	return http.ListenAndServe(`:8080`, mux)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.Header.Add("Content-Type", "text/plain")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	shortenStr := "hevfyegruf"
	response := fmt.Sprintf("http://localhost:8080/%s", shortenStr)
	urlStore[shortenStr] = originalURL

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]

	mu.Lock()
	defer mu.Unlock()

	originalURL, exists := urlStore[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
