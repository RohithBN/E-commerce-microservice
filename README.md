

# Start Redis
docker run -d --name redis -p 6379:6379 redis

# Start Kafka and Zookeeper
docker run -d --name zookeeper -p 2181:2181 wurstmeister/zookeeper
docker run -d --name kafka -p 9092:9092 \
  -e KAFKA_ADVERTISED_HOST_NAME=localhost \
  -e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 \
  wurstmeister/kafka

# Start Prometheus
docker run -d --name prometheus \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus

# Start Grafana
docker run -d --name grafana -p 3000:3000 grafana/grafana



# Start each service in separate terminal windows

# Start Auth Service
cd auth-service && go run main.go

# Start Product Service
cd product-service && go run main.go

# Start Cart Service
cd cart-service && go run main.go

# Start Order Service
cd order-service && go run main.go

# Start Gateway
cd gateway && go run main.go


#Authentication Testing
# Register a new user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User 2","email":"test2@example.com","password":"password123"}'

# Login to get a JWT token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Store the token in an environment variable
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."



#Product Service Testing
# Add a new product
curl -X POST http://localhost:8080/api/add-product \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "This is a test product",
    "price": 99.99,
    "stock": 100,
    "category": "Electronics"
  }'

# Get all products
curl -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN"

# Get a specific product
# Replace PRODUCT_ID with an actual ID from the previous response
export PRODUCT_ID="..."
curl -X GET "http://localhost:8080/api/products/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN"

# Update a product
curl -X PUT "http://localhost:8080/api/update-product/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "description": "This is an updated product",
    "price": 149.99,
    "stock": 50,
    "category": "Electronics"
  }'

# Delete a product
curl -X DELETE "http://localhost:8080/api/delete-product/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN"



#Cart service testing
# Add a product to cart
curl -X POST "http://localhost:8080/api/cart/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '2'

# Get cart contents
curl -X GET http://localhost:8080/api/cart \
  -H "Authorization: Bearer $TOKEN"

# Remove a product from cart
curl -X DELETE "http://localhost:8080/api/cart/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN"


  
# Order Service Testing
# Add product to cart first
curl -X POST "http://localhost:8080/api/cart/$PRODUCT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '1'

# Request OTP for 2FA
curl -X POST http://localhost:8080/api/orders/send-otp \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

# Verify OTP
curl -X POST http://localhost:8080/api/orders/verify-otp \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"otp":"123456"}'

# Create a new order from cart
curl -X POST http://localhost:8080/api/create-order \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shippingAddress": "123 Test Street, City, 12345",
    "paymentMethod": "credit_card"
  }'

# Save the order ID from the response
export ORDER_ID="..."

# Process payment for the order
curl -X POST http://localhost:8080/api/orders/payment \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "'$ORDER_ID'",
    "paymentMethod": "credit_card",
    "cardNumber": "4242424242424242",
    "expiryMonth": "12",
    "expiryYear": "2025",
    "cvv": "123"
  }'

# Update order status
curl -X PUT "http://localhost:8080/api/orders/$ORDER_ID/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "shipped"
  }'

# Get all orders
curl -X GET http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN"