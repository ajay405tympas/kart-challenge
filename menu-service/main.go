package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://orderfoodonline.deno.dev/api"

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// CORS middleware
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key")
}

// GET /product
func listProducts(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.URL.Path != "/product" {
		http.NotFound(w, r)
		return
	}

	log.Println("Request received: /product")

	resp, err := client.Get(baseURL + "/product")
	if err != nil {
		log.Println("Error calling upstream:", err)
		http.Error(w, "Upstream service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println("Error writing response:", err)
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

	log.Println("Request received: /product/", id)

	resp, err := client.Get(baseURL + "/product/" + id)
	if err != nil {
		log.Println("Error calling upstream:", err)
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

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func main() {
	http.HandleFunc("/product", listProducts)
	http.HandleFunc("/product/", getProduct)

	log.Println("Menu service running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}