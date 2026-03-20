package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Req struct {
	Coupon string  `json:"coupon"`
	Amount float64 `json:"amount"`
}

type Resp struct {
	FinalAmount float64 `json:"finalAmount"`
}

type CouponsResp struct {
	Coupons []string `json:"coupons"`
}

// ✅ Get env with fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func apply(code string, amt float64) float64 {
	switch code {
	case "DISCOUNT10":
		return amt * 0.9
	case "FLAT20":
		return amt - 20
	case "MAR10":
		return amt * 0.9
	case "FIRSTTIME15":
		return amt * 0.85
	default:
		return amt
	}
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

	port := getEnv("PORT", "8085")

	http.HandleFunc("/coupons", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {

		log.Println("📥 Fetching available coupons")

		coupons := []string{"FLAT20", "MAR10", "DISCOUNT10", "FIRSTTIME15"}

		resp := CouponsResp{
			Coupons: coupons,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))

	http.HandleFunc("/paywithcoupon", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {

		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("❌ Invalid request:", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		log.Printf("🎯 Applying coupon: %s on amount %.2f\n", req.Coupon, req.Amount)

		finalAmount := apply(req.Coupon, req.Amount)

		res := Resp{
			FinalAmount: finalAmount,
		}

		log.Printf("💸 Final amount after discount: %.2f\n", finalAmount)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}))

	log.Println("🚀 Coupon service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
