package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	// mu       sync.Mutex
	mu = sync.Mutex{}
)

func main() {
	// mux := http.NewServeMux()

	// mux.HandleFunc("/", postHandler)
	// mux.HandleFunc("/{id}", getHandler)
	// fmt.Println("Server started at http://localhost:8080")

	// err := http.ListenAndServe(`:8080`, mux)
	// if err != nil {
	// 	panic(err)
	// }
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", postHandler)
	// mux.HandleFunc("/{id}", getHandler) // try to pass an argument
	mux.HandleFunc("/hevfyegruf", getHandler) // try to pass an argument
	fmt.Println("Server started at http://localhost:8080")

	return http.ListenAndServe(`:8080`, mux)
}

// func shortenURL() string {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	seed := rand.NewSource(time.Now().UnixNano())
// 	random := rand.New(seed)

// 	result := make([]byte, 8)
// 	for i := range result {
// 		result[i] = charset[random.Intn(len(charset))]
// 	}

// 	return string(result)
// }

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.Header.Add("Content-Type", "text/plain") // ???

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest) // TODO
		return
	}
	defer r.Body.Close()

	originalURL := string(body) // https://practicum.yandex.ru/
	if originalURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest) // TODO
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// if _, exists := urlStore[originalURL]; !exists {
	// 	urlStore[originalURL] = shortenURL()
	// }
	// response := "http://localhost:8080/" + urlStore[originalURL]
	// w.Header().Set("Content-Type", "text/plain")
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte(response))

	shortenStr := "hevfyegruf"
	response := fmt.Sprintf("http://localhost:8080/%s", shortenStr) // "http://localhost:8080/123"
	urlStore[shortenStr] = originalURL                              // { 123 => https://practicum.yandex.ru/ }

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]

	mu.Lock()
	defer mu.Unlock()
	// for originalURL, shortenedID := range urlStore {
	// 	if shortenedID == id {

	// 		w.WriteHeader(http.StatusTemporaryRedirect)
	// 		w.Header().Set("Location", originalURL)
	// 		return
	// 	}
	// }

	// http.Error(w, "Shortened URL not found", http.StatusNotFound)
	fmt.Println(urlStore)
	originalURL, exists := urlStore[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
