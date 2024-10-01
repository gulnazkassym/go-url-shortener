package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	urlStore = make(map[string]string)
	mu       sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", postHandler)
	mux.HandleFunc("/{id}", getHandler)
	fmt.Println("Server started at http://localhost:8080")

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func shortenURL() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest) // TODO
		return
	}
	defer r.Body.Close()

	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest) // TODO
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := urlStore[originalURL]; !exists {
		urlStore[originalURL] = shortenURL()
	}
	response := "http://localhost:8080/" + urlStore[originalURL]
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]

	mu.Lock()
	defer mu.Unlock()
	for originalURL, shortenedID := range urlStore {
		if shortenedID == id {
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Header().Set("Location", originalURL)
			return
		}
	}

	http.Error(w, "Shortened URL not found", http.StatusNotFound)
}
