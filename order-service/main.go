package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type Req struct {
	Items  []Item `json:"items"`
	Coupon string `json:"couponCode"`
}

// ✅ Get env with fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// ✅ CORS middleware (env-driven)
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		origin := getEnv("ALLOWED_ORIGIN", "http://localhost:3000")

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {

	port := getEnv("PORT", "8083")
	couponServiceURL := getEnv("COUPON_SERVICE_URL", "http://localhost:8085")

	http.HandleFunc("/order", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {

		log.Println("➡️ Incoming /order request")

		// ✅ Decode request safely
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("❌ Decode error:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// ✅ Validate request
		if len(req.Items) == 0 {
			http.Error(w, "No items in order", http.StatusBadRequest)
			return
		}

		log.Printf("📦 Request: %+v\n", req)

		// ✅ Calculate total
		total := 0.0
		for _, i := range req.Items {
			switch i.ProductID {
			case "1":
				total += 100 * float64(i.Quantity)
			case "2":
				total += 200 * float64(i.Quantity)
			default:
				log.Println("⚠️ Unknown product ID:", i.ProductID)
			}
		}

		log.Println("💰 Calculated total:", total)

		// ✅ Prepare request to coupon service
		body, err := json.Marshal(map[string]interface{}{
			"amount": total,
			"coupon": req.Coupon,
		})
		if err != nil {
			log.Println("❌ Marshal error:", err)
			http.Error(w, "Failed to prepare request", http.StatusInternalServerError)
			return
		}

		// ✅ Call coupon service (env-driven URL)
		url := couponServiceURL + "/paywithcoupon"
		log.Println("🔄 Calling Coupon Service:", url)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Println("❌ Coupon service not reachable:", err)
			http.Error(w, "Coupon service not reachable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// ✅ Read response body
		respBody, _ := io.ReadAll(resp.Body)
		log.Println("📩 Coupon service response:", string(respBody))

		// ✅ Handle non-200 responses
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Coupon service error: "+string(respBody), http.StatusBadGateway)
			return
		}

		// ✅ Decode response
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			log.Println("❌ Decode response error:", err)
			http.Error(w, "Invalid response from coupon service", http.StatusInternalServerError)
			return
		}

		// ✅ Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))

	log.Println("🚀 Order service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
