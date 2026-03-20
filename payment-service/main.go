package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Req struct {
	Amount float64 `json:"amount"`
	Coupon string  `json:"coupon"`
}

type Resp struct {
	Status           string  `json:"status"`
	OriginalAmount   float64 `json:"originalAmount"`
	DiscountedAmount float64 `json:"discountedAmount"`
	FinalAmount      float64 `json:"finalAmount"`
	Message          string  `json:"message"`
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

	port := getEnv("PORT", "8084")
	couponServiceURL := getEnv("COUPON_SERVICE_URL", "http://localhost:8085")

	log.Println("💳 Payment service starting on port:", port)

	http.HandleFunc("/pay", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {

		log.Println("➡️ Incoming /pay request")

		// ✅ Decode request
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("❌ Decode error:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("📥 Request: Amount=%.2f, Coupon=%s\n", req.Amount, req.Coupon)

		// ✅ Default discounted amount
		discountedAmount := req.Amount

		// ✅ Call coupon service
		if req.Coupon != "" {
			body, _ := json.Marshal(map[string]interface{}{
				"coupon": req.Coupon,
				"amount": req.Amount,
			})

			url := couponServiceURL + "/paywithcoupon"
			log.Println("🔄 Calling Coupon Service:", url)

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Println("⚠️ Coupon service not reachable:", err)
			} else {
				defer resp.Body.Close()

				var data map[string]float64
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					log.Println("⚠️ Decode coupon response error:", err)
				} else {
					if val, ok := data["finalAmount"]; ok {
						discountedAmount = val
						log.Printf("💸 Discounted amount: %.2f\n", discountedAmount)
					}
				}
			}
		}

		// ✅ Simulate payment
		log.Printf("💳 Processing payment for: %.2f\n", discountedAmount)

		finalAmount := 0.0

		log.Println("✅ Payment successful")

		// ✅ Response
		response := Resp{
			Status:           "SUCCESS",
			OriginalAmount:   req.Amount,
			DiscountedAmount: discountedAmount,
			FinalAmount:      finalAmount,
			Message:          "Order completed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
