package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var (
	urlStore = make(map[string]string)
	mu       = sync.Mutex{}
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	r.Post("/", postHandler)
	r.Get("/hevfyegruf", getHandler)

	fmt.Println("Server started at", flagRunAddr)

	return http.ListenAndServe(flagRunAddr, r)
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

	// originalURL, exists := urlStore[id]
	// if !exists {
	// 	http.NotFound(w, r)
	// 	return
	// }
	originalURL := urlStore[id]

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
