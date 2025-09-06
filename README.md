# Food Ordering System
Food Ordering System built with Golang

## Features
- Add menu items with name, price, availability
- Place an order with multiple items
- Generate a bill with tax
- Update menu item availability

## Design

## Entities

### Menu
- menu will just be a list of items

### MenuItem
- item contains: name, price, availability

### Order
- order is a struct with a list of menu-items x quantity and customer-id (optional)

### Invoice/Bill
- Invoice will contian the order info (list of items with price) with all the taxes, payment status (done or not)

### Users
- Customers: who place orders
- Restaurant Owner: who adds items or manage menu

### Restaurants
- Restuarant is a entity, it can have it's own menu items, customers can place orders in a restaurant

## Future Enhancements
### Taxes
- List of taxes for differnet restaurants restaurant, eg, gst on food (cgst & sgst), restaurant service charges/tax, etc...
- Tax will contain: type (fixed charge/percentage), name, description, tax (float)

## APIs

### Authentication
- `POST /api/auth/login` 
- `POST /api/auth/logout`

### Users
- `POST /api/users`
- `GET /api/users/{id}`
<!-- - `PUT /api/users/{id}` -->
<!-- - `DELETE /api/users/{id}` -->

### Restaurants
- `GET /api/restaurants`
- `POST /api/restaurants`
- `GET /api/restaurants/{id}`
- `GET /api/restaurants/{id}/items`
- `POST /api/restaurants/{id}/items`
<!-- - `PUT /api/restaurants/{id}` -->
<!-- - `DELETE /api/restaurants/{id}` -->

### Menu Items
- `GET /api/items/{id}`
- `PATCH /api/items/{id}` (availability)
<!-- - `PUT /api/items/{id}` -->
<!-- - `DELETE /api/items/{id}` -->

### Orders
- `GET /api/orders?user_id=<id>`
- `POST /api/orders`
- `GET /api/orders/{id}`
- `POST /api/orders/{id}/items`
- `GET /api/orders/{id}/invoice`
<!-- - `PATCH /api/orders/{id}` -->
<!-- - `DELETE /api/orders/{id}` -->
