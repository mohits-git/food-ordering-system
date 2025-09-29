# Food Ordering System
Food Ordering System built with Golang

## Project Setup

This project contains http APIs for the food ordering system and a cli client for basic interaction with flow with apis

- Copy envs
```bash
cp .env.example .env
```

- Install dependencies
```bash
go mod tidy
```

- Run the server
```bash
go run ./cmd/api
```

- Run the cli client
```bash
go run ./cmd/client
```

- Run Tests
```bash
go test -cover -coverprofile=cover.out ./internal/...
```

- With tparse
```bash
go install github.com/mfridman/tparse@latest
```
```bash
go test -cover -coverprofile=cover.out -json ./internal/... | tparse -all
```

## Features
- Add menu items with name, price, availability
- Place an order with multiple items
- Generate a bill with tax
- Update menu item availability

## Entities

### Menu
- menu will just be a list of items ### MenuItem item contains: name, price, availability

### Order
- order is a struct with a list of menu-items x quantity and customer-id (optional)

### Invoice/Bill
- Invoice will contian the order info (list of items with price) with all the taxes, payment status (done or not)

### Users
- Customers: who place orders
- Restaurant Owner: who adds items or manage menu

### Restaurants
- Restuarant is a entity, it can have it's own menu items, customers can place orders in a restaurant

## APIs

### Authentication
- `POST /api/auth/login` 
- `POST /api/auth/logout` (authenticated)

### Users
- `POST /api/users`
- `GET /api/users/{id}`
<!-- - `PUT /api/users/{id}` -->
<!-- - `DELETE /api/users/{id}` -->

### Restaurants
- `GET /api/restaurants`
- `POST /api/restaurants` (authenticated)
<!-- - `GET /api/restaurants/{id}` -->
<!-- - `PUT /api/restaurants/{id}` -->
<!-- - `DELETE /api/restaurants/{id}` -->

### Menu Items
- `GET /api/restaurants/{id}/items`
- `POST /api/restaurants/{id}/items` (authenticated)
- `PATCH /api/items/{id}` (availability) (authenticated)
<!-- - `GET /api/items/{id}` -->
<!-- - `PUT /api/items/{id}` -->
<!-- - `DELETE /api/items/{id}` -->

### Orders
- `POST /api/orders` (authenticated)
- `POST /api/orders/{id}/items` (authenticated)
- `GET /api/orders/{id}` (authenticated)
<!-- - `GET /api/orders?user_id=<id>` -->
<!-- - `PATCH /api/orders/{id}` -->
<!-- - `DELETE /api/orders/{id}` -->

### Invoice
- `POST /api/orders/{id}/invoices` (authenticated)
- `POST /api/invoices/{id}/pay` (authenticated)
- `GET /api/invoices/{id}` (authenticated)
