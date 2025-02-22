# Order Matching API

## Overview
This project is an Order Matching API built with Go and Gin. It allows users to place buy and sell orders, retrieve the order book, and view existing orders. 

## Features
- Place buy and sell orders
- Retrieve order book
- Get a list of existing orders
- Concurrency handling with mutex locks
- Swagger API documentation

## Installation

### Steps
1. Clone the repository:
   ```sh
   git clone https://github.com/mta9896/order-matching.git
   cd order-matching
   ```
2. Run the application using docker-compose:
   ```sh
   docker-compose up --build
   ```

## API Documentation
Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```

To generate or update Swagger documentation, run:
```sh
swag init
```

## API Endpoints
### 1. Place Order
**POST /api/orders**
- Places a buy or sell order.

### 2. Get Order Book
**GET /api/orderbook?limit=10**
- Retrieves the current state of the order book.

### 3. Get Orders List
**GET /api/orders?page=1&page_size=10**
- Returns a paginated list of orders.

## Concurrency Handling
To prevent race conditions when placing orders, a mutex lock is used in `CreateOrder` to ensure safe access to shared resources. This prevents duplicate order processing and ensures thread safety.

## Author
Marzieh Tajik - [GitHub Profile](https://github.com/mta9896)

