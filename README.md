# Ecommerce API

Ecommerce API developed in Golang responsable for managing customers, products and orders, allowing a wide set of actions such as fetching, creating, deleting, paying and cancelling customer orders, managing product stocks, prices and customer data.

Based mainly on `net/http` and `chi`, it also includes `pgxpool` for connecting to Postgres, `godotenv` for dealing with environment variables, `bcrypt` for password hashing, among many other libraries.

For more detailed specifications, check the requirements doc here: [Challenge Document](challenge.md)

## Stack

*   **Golang**
*   **PostgreSQL** (integrated using `pgx`)
*   **Docker**
*   **Migrations**

## Endpoints

### Customers

#### Get all customers

```http
  GET /customer
```

#### Create customer

```http
  POST /customer
```

#### Update customer

```http
  PUT /customer/${customer_id}
```

#### Delete customer

```http
  DELETE /customer/${customer_id}
```

### Products

#### Get all products

```http
  GET /product
```

#### Create product

```http
  POST /product
```

#### Update product

```http
  PUT /product/${product_id}
```

#### Delete product

```http
  DELETE /product/${product_id}
```

### Orders

#### Get all orders by customer

```http
  GET /order/all/${customer_id}
```

#### Get order by ID

```http
  GET /order/${order_id}
```

#### Create order

```http
  POST /order
```

#### Cancel order

```http
  POST /order/cancel/${order_id}
```

#### Pay order

```http
  POST /order/pay/${order_id}
```

## How to run this API:

```bash
  docker compose up -d
  migrate -path db/migrations -database $database_link up
  go run cmd/api/main.go
```
