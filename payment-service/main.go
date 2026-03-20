package main

import (
	"bytes"
	"encoding/json"
	"net/http"
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

	http.HandleFunc("/pay", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {

		// ✅ Decode request
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// ✅ Default discounted amount
		discountedAmount := req.Amount

		// ✅ Call coupon service (if coupon exists)
		if req.Coupon != "" {
			body, _ := json.Marshal(map[string]interface{}{
				"code":   req.Coupon,
				"amount": req.Amount,
			})

			resp, err := http.Post("http://localhost:8085/paywithcoupon", "application/json", bytes.NewBuffer(body))
			if err == nil && resp != nil {
				defer resp.Body.Close()

				var data map[string]float64
				if json.NewDecoder(resp.Body).Decode(&data) == nil {
					if val, ok := data["finalAmount"]; ok {
						discountedAmount = val
					}
				}
			}
			// graceful fallback if coupon fails
		}

		// ✅ Simulate payment processing using discounted amount
		processedAmount := discountedAmount

		// (used to avoid unused variable + realistic flow)
		_ = processedAmount

		// ✅ Payment success → balance cleared
		finalAmount := 0.0

		// ✅ Response to client
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

	http.ListenAndServe(":8084", nil)
}
