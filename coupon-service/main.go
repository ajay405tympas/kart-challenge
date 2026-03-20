package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Req struct {
	Coupon string  `json:"coupon"` // ✅ FIXED
	Amount float64 `json:"amount"`
}

type Resp struct {
	FinalAmount float64 `json:"finalAmount"`
}

type CouponsResp struct {
	Coupons []string `json:"coupons"`
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

// ✅ CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
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

	http.HandleFunc("/coupons", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {

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
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		log.Println("🎯 Received coupon:", req.Coupon)

		finalAmount := apply(req.Coupon, req.Amount)

		res := Resp{
			FinalAmount: finalAmount,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}))

	log.Println("Coupon service running on :8085")
	http.ListenAndServe(":8085", nil)
}
