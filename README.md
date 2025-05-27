Here's a single `README.md` file that combines both **Order Service** and **Payment Service** endpoints, deployment instructions, and usage in **one unified file**:

---

````markdown
# 🧾 Order & Payment Microservices (Go + Cloud Run + Pub/Sub)

This project includes two microservices:

- **Order Service**: Creates orders, publishes to Pub/Sub
- **Payment Service**: Subscribes to Pub/Sub, processes payments

---

## 🌐 Base URLs (After Deployment)

Replace with actual Cloud Run URLs:

- `ORDER_SERVICE_URL=https://order-service-357930214451.us-central1.run.app`
- `PAYMENT_SERVICE_url=https://payment-service-357930214451.us-central1.run.app`

---

## 🔐 Authentication

All endpoints require a JWT token.

### ✅ Generate JWT Token

Use this Go snippet:

```go
package main

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	secret := "your_jwt_secret" // same as JWT_SECRET in .env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signedToken, _ := token.SignedString([]byte(secret))
	fmt.Println("JWT Token:", signedToken)
}
````

---

## 📦 .env Example

```env
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=postgres
CLOUDSQL_CONNECTION_NAME=micro-services-tech-ready:us-central1:your-instance

PROJECT_ID=micro-services-tech-ready
TOPIC_ID=order-events
SUBSCRIPTION_ID=payment-sub

JWT_SECRET=your_jwt_secret
```

---

## 📚 API Endpoints (Order + Payment)

> 💡 All requests must include header:
> `Authorization: Bearer <JWT_TOKEN>`

---

### 🆕 Create Order

**POST `/orders`**

Request:

```json
{
  "user_id": "123",
  "items": [
    {
      "product_id": "p1",
      "product_name": "Almonds",
      "quantity": 2,
      "price": 100
    },
    {
      "product_id": "p2",
      "product_name": "Cashews",
      "quantity": 1,
      "price": 150
    }
  ]
}
```

Response:

```json
{
  "order_id": "uuid",
  "total_amount": 350
}
```

---

### 📄 Get Order By ID

**GET `/orders/:id`**

Response:

```json
{
  "id": "uuid",
  "user_id": "123",
  "items": [...],
  "status": "CREATED",
  "total_amount": 350
}
```

---

### ✏️ Update Order Status

**PATCH `/orders/:id/status`**

Request:

```json
{
  "status": "SHIPPED"
}
```

Response:

```json
{ "message": "Order status updated" }
```

---

### 💳 Payment Service

The Payment Service **has no public HTTP endpoints**. It listens to Pub/Sub and:

* Saves a new payment record when an order is received
* Logs the result

Check logs via:

```bash
gcloud logging read "resource.labels.service_name=payment-service" --limit=10
```

---

## ☁️ Deployment Guide

### 🐳 Docker Build & Push

```bash
# Order Service
docker build -t us-central1-docker.pkg.dev/YOUR_PROJECT/my-docker-repo/order-service .
docker push us-central1-docker.pkg.dev/YOUR_PROJECT/my-docker-repo/order-service

# Payment Service
docker build -t us-central1-docker.pkg.dev/YOUR_PROJECT/my-docker-repo/payment-service .
docker push us-central1-docker.pkg.dev/YOUR_PROJECT/my-docker-repo/payment-service
```

---

### 🚀 Deploy to Cloud Run

```bash
# Order Service (public)
gcloud run deploy order-service \
  --image=us-central1-docker.pkg.dev/YOUR_PROJECT/my-docker-repo/order-service \
  --region=us-central1 \
  --allow-unauthenticated \
  --port=8080

# Payment Service (private)
gcloud run deploy payment-service \
  --image=us-central1-docker.pkg.dev/YOUR_PROJECT/my-docker-repo/payment-service \
  --region=us-central1 \
  --no-allow-unauthenticated \
  --memory=512Mi \
  --timeout=300s
```

---

### 🔗 Setup Pub/Sub

```bash
# Topic
gcloud pubsub topics create order-events

# Subscription
gcloud pubsub subscriptions create payment-sub \
  --topic=order-events
```

---

## 🧪 Test Workflow

1. Generate JWT Token
2. Call `POST /orders`
3. Order service:

   * Saves order to DB
   * Publishes event to Pub/Sub
4. Payment service:

   * Subscribes to Pub/Sub
   * Inserts payment record
5. Check DB or Cloud Logs

---

## 🐘 Local Dev (with Cloud SQL)

Run Cloud SQL proxy:

```bash
./cloud_sql_proxy -dir=/cloudsql -instances="micro-services-tech-ready:us-central1:your-instance"
```

Use this in `.env`:

```env
DB_HOST=/cloudsql/micro-services-tech-ready:us-central1:your-instance
```

---

## 🧯 Troubleshooting

| Issue                       | Fix                                           |
| --------------------------- | --------------------------------------------- |
| `no such file or directory` | Run Cloud SQL proxy                           |
| `Pub/Sub not working`       | Check topic/subscription names and IAM roles  |
| `DB connect failed`         | Verify env vars and socket path               |
| `payment-service exits`     | It's fine, it will restart on Pub/Sub trigger |

---

## ✅ Permissions

| Role                 | Required For       |
| -------------------- | ------------------ |
| `Cloud Run Invoker`  | Call Order Service |
| `Pub/Sub Publisher`  | Order Service      |
| `Pub/Sub Subscriber` | Payment Service    |
| `Cloud SQL Client`   | Both services      |

---

## 🏁 Done!

You now have:

* JWT-secured Order API
* Event-driven Payment processing
* Cloud SQL integration
* Cloud Run deployment

🚀 Ready for production.

```

Let me know if you want:
- This formatted as a downloadable file  
- Sample `.env`, `.sql` schema  
- Postman collection  
- cURL test scripts
```
