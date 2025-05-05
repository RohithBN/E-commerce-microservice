

### Prerequisites

Ensure you have Docker, Go, and `curl` installed.

### 🐳 Start Dependencies

```bash
# Redis
docker run -d --name redis -p 6379:6379 redis

# Zookeeper & Kafka
docker run -d --name zookeeper -p 2181:2181 wurstmeister/zookeeper
docker run -d --name kafka -p 9092:9092 \
  -e KAFKA_ADVERTISED_HOST_NAME=localhost \
  -e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 \
  wurstmeister/kafka

# Prometheus
docker run -d --name prometheus -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus

# Grafana
docker run -d --name grafana -p 3000:3000 grafana/grafana
```

### 🧩 Start Go Services

Open separate terminals and run each service:

```bash
cd auth-service && go run main.go
cd product-service && go run main.go
cd cart-service && go run main.go
cd order-service && go run main.go
cd gateway && go run main.go
```

---

## 🔐 Authentication

```bash
# Register
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Save token
export TOKEN="your_jwt_token_here"
```

---

## 🛍️ Product APIs

```bash
# Add Product
curl -X POST http://localhost:8080/api/add-product \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","description":"A product","price":99.99,"stock":100,"category":"Electronics"}'

# Get Products
curl -X GET http://localhost:8080/api/products -H "Authorization: Bearer $TOKEN"
```

---

## 🛒 Cart APIs

```bash
# Add to Cart
curl -X POST "http://localhost:8080/api/cart/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN" -d '2'

# Get Cart
curl -X GET http://localhost:8080/api/cart -H "Authorization: Bearer $TOKEN"
```

---

## 📦 Order & 2FA Flow

```bash
# Send OTP
curl -X POST http://localhost:8080/api/orders/send-otp \
  -H "Authorization: Bearer $TOKEN"

# Verify OTP
curl -X POST http://localhost:8080/api/orders/verify-otp \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"otp":"123456"}'

# Create Order
curl -X POST http://localhost:8080/api/create-order \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"shippingAddress":"123 Main St","paymentMethod":"credit_card"}'

# Save order ID
export ORDER_ID="your_order_id_here"

# Payment
curl -X POST http://localhost:8080/api/orders/payment \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"orderId":"'$ORDER_ID'","paymentMethod":"credit_card","cardNumber":"4242...","expiryMonth":"12","expiryYear":"2025","cvv":"123"}'
```

---

