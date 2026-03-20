# 🍔 Food Ordering Microservices Project

## 📌 Overview

This project is a **microservices-based food ordering system** built
using: - React (UI Service) - Go (Backend Services) - Docker &
Kubernetes (Deployment Ready)

------------------------------------------------------------------------

## 🧩 Services Architecture

    UI (React - localhost:3000)
       ↓
    Order Service (8083)
       ↓
    Payment Service (8084)
       ↓
    Coupon Service (8085)
       ↓
    Menu Service (8080)

------------------------------------------------------------------------

## 🚀 Services Description

### 1. UI Service (React)

-   Displays products
-   Allows adding items to cart
-   Applies coupons
-   Calls `/order` API

------------------------------------------------------------------------

### 2. Menu Service

-   Fetches products from upstream API
-   Endpoints:
    -   GET /product
    -   GET /product/{id}

------------------------------------------------------------------------

### 3. Order Service

-   Accepts cart + coupon
-   Calculates total
-   Calls Coupon Service
-   Returns final amount

------------------------------------------------------------------------

### 4. Payment Service

-   Processes payment
-   Applies coupon via Coupon Service
-   Returns order success

------------------------------------------------------------------------

### 5. Coupon Service

-   Lists available coupons
-   Applies discount logic

------------------------------------------------------------------------

## ⚙️ Environment Variables

Each service uses `.env` for config:

### Example

    PORT=8083
    ALLOWED_ORIGIN=http://localhost:3000
    COUPON_SERVICE_URL=http://coupon-service:8085

------------------------------------------------------------------------

## 🐳 Docker Support

Each service can be containerized using Docker.

Use service names for communication:

    http://coupon-service:8085
    http://payment-service:8084

------------------------------------------------------------------------

## ☸️ Kubernetes Deployment

### Example Service DNS

    http://coupon-service:8085

### Deployment ENV Example

    env:
    - name: COUPON_SERVICE_URL
      value: "http://coupon-service:8085"

------------------------------------------------------------------------

## 🔥 Key Features

-   ✅ CORS enabled
-   ✅ Environment-based config
-   ✅ Service-to-service communication
-   ✅ Logging enabled
-   ✅ Error handling
-   ✅ Kubernetes-ready

------------------------------------------------------------------------

## 🧪 Running Locally

1.  Start services:

```{=html}
<!-- -->
```
    go run main.go

2.  Start UI:

```{=html}
<!-- -->
```
    npm start

------------------------------------------------------------------------

## 📡 API Flow

1.  UI → Order Service
2.  Order → Payment Service
3.  Payment → Coupon Service
4.  Response back to UI

------------------------------------------------------------------------

## 🧠 Best Practices Used

-   12-Factor App principles
-   Environment-based configuration
-   Microservices separation
-   Graceful error handling

------------------------------------------------------------------------

## 🚀 Future Enhancements

-   Add database (orders, users)
-   Add authentication (JWT)
-   Add API Gateway / Ingress
-   Add observability (Prometheus, Grafana)

------------------------------------------------------------------------

## 👨‍💻 Author

Built for learning microservices, Go, Docker & Kubernetes.
