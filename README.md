# 🍽️ Restaurant API

A secure, multi-tenant REST API built in Go using the Fiber framework. It integrates with Square POS and supports creating orders, retrieving orders, and processing payments. Also includes robust health checks and metrics.

---

## 📦 Local Development

### ✅ Prerequisites

- Go 1.21+
- PostgreSQL running locally on default port (`5432`)

### 📄 Environment Setup

Create a `.env` file in the root directory:

```env
DSN=postgres://postgres@localhost:5432/restaurant_api?sslmode=disable
PORT=3003
```

#### 🚀 Run the Server

```bash
go run ./cmd/api
```

### 🐳 Docker Setup

A Docker-based setup is provided, which includes PostgreSQL and all dependencies.
🏁 Start with Docker Compose

```bash
docker compose up --build
```

This will:

    Start a PostgreSQL database (restaurant-db)

    Start the API server (restaurant-api)

    Expose the API on http://localhost:3003

### 🩺 Health Checks

Health checks are included using fiber/middleware/healthcheck.

app.Use(healthcheck.New())

## 🔍 Endpoints

Liveness Probe

    Endpoint: /livez

        Behavior: Returns 200 OK when the server is running.

    Readiness Probe

        Endpoint: /readyz

        Behavior: Returns 200 OK when the application is ready to handle requests.

If the application is not ready, it returns 503 Service Unavailable.
📊 Metrics

The /metrics endpoint provides application performance metrics using fiber/middleware/monitor.

```go
fiber.Get("/metrics", monitor.New(monitor.Config{
    Title: "Square POS API Metrics Page",
}))
```

Access at: http://localhost:3003/metrics
### 🔐 Authenticated Routes

All routes are grouped under the /v1 prefix and require an Authorization token in the request header.
🔒 Header

```
Authorization: <YOUR_TOKEN>
Content-Type: application/json
```

| Method | Endpoint                          | Description                   |
|--------|-----------------------------------|-------------------------------|
| POST   | `/v1/orders`                      | Create a new order            |
| GET    | `/v1/orders/:id`                  | Get order by ID               |
| GET    | `/v1/orders/table/:tableNumber`   | Get orders for a table        |
| POST   | `/v1/orders/:orderId/pay`         | Process payment for an order  |

🧪 Sample Requests

All requests use port 3003.

🔸 Create Order

```bash
curl -X POST 'http://localhost:3003/v1/orders' \
  --header 'Authorization: EAAAl7Y-od7IFd0hK3kB4loclod4MVyxd9ol2VGlLN1J1WH1-ymXWz8PrbxXYXgq' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "tableNumber": "121",
    "items": [
      {
        "name": "Burger",
        "quantity": 2,
        "unitPrice": 1200
      },
      {
        "name": "Fries",
        "quantity": 1,
        "unitPrice": 500
      }
    ]
  }'
```

🔸 Process Payment

```bash
curl -X POST 'http://localhost:3003/v1/orders/SB9D03sB4A5yM4YS1FksERNNXPTZY/pay' \
  --header 'Authorization: EAAAl7Y-od7IFd0hK3kB4loclod4MVyxd9ol2VGlLN1J1WH1-ymXWz8PrbxXYXgq' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "billAmount": 2900.00,
    "tipAmount": 0.00,
    "paymentId": "G6M56S"
  }'
```

🔸 Get Orders by Table Number

```bash
curl -X GET 'http://localhost:3003/v1/orders/table/123' \
  --header 'Authorization: EAAAl7Y-od7IFd0hK3kB4loclod4MVyxd9ol2VGlLN1J1WH1-ymXWz8PrbxXYXgq'
```

🔸 Get Order by ID

```bash
curl -X GET 'http://localhost:3003/v1/orders/e1k7WQyRMYFIq8WUNGLNYlVI8y8YY' \
  --header 'Authorization: EAAAl7Y-od7IFd0hK3kB4loclod4MVyxd9ol2VGlLN1J1WH# Markdown syntax guide
```