package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// ✅ Load env variables with fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// CORS middleware
func enableCORS(w http.ResponseWriter, r *http.Request) {
	origin := getEnv("ALLOWED_ORIGIN", "http://localhost:3000")

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key")
}

// GET /product
func listProducts(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.URL.Path != "/product" {
		http.NotFound(w, r)
		return
	}

	baseURL := getEnv("BASE_URL", "https://orderfoodonline.deno.dev/api")

	log.Println("📥 Request: /product")

	resp, err := client.Get(baseURL + "/product")
	if err != nil {
		log.Println("❌ Upstream error:", err)
		http.Error(w, "Upstream service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("❌ Response write error:", err)
	}
}

// GET /product/{id}
func getProduct(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/product/")
	if id == "" {
		http.Error(w, "Missing productId", http.StatusBadRequest)
		return
	}

	baseURL := getEnv("BASE_URL", "https://orderfoodonline.deno.dev/api")

	log.Println("📥 Request: /product/", id)

	resp, err := client.Get(baseURL + "/product/" + id)
	if err != nil {
		log.Println("❌ Upstream error:", err)
		http.Error(w, "Upstream service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("❌ Response write error:", err)
	}
}

func main() {
	port := getEnv("PORT", "8080")

	http.HandleFunc("/product", listProducts)
	http.HandleFunc("/product/", getProduct)

	log.Println("🚀 Menu service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
