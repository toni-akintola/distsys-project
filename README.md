South Bend Bets is a mock trading platform where users can create accounts, view stock data, and place trade orders. It consists of two main servers:
1. **Market Server**: Provides stock data and handles price updates.
2. **Executor Server**: Manages user accounts and processes trade orders.

---

## Table of Contents
- [Installation](#installation)
- [How to Run](#how-to-run)
- [API Usage](#api-usage)
  - [Account Creation](#account-creation)
  - [Fetching Account Info](#fetching-account-info)
  - [Viewing Stock Data](#viewing-stock-data)
  - [Placing Orders](#placing-orders)

---

## Installation (using any of the nd.edu student machines)

1. **Clone the Repository**:
```bash
   git clone <repository-url>
   cd <repository-folder>
```

2. **Version check**:
    Ensure that Go 1.20.x is installed on your machine.

## How to Run
1. **Start the market server & executor server instances**:
    `go run main.go market-server.go executor-server.go`

2. **Run `client_setup.py` script to retrieve the URL of the executor from the name server.**
---

## API Usage
### Account Creation
`Endpoint: /create-account`
`Method: POST`
`Request Body:`
```json
{
  "username": "example_user"
}
```

`Response:`
```json
{
  "username": "example_user",
  "balance": 10000,
  "positions": []
}
```
---

### Fetching Account Info

**Endpoint**: `/account/{username}`  
**Method**: `GET`  
**Response**:
```json

{"username": "example_user", "balance": 9500, "positions": [{"order": {"quantity": 10, "ticker": "AAPL", "username": "example_user"}, "price": 150}]}
```
---

### Viewing Stock Data

#### Single Stock

**Endpoint**: `/single-stock/{ticker}`  
**Method**: `GET`  
**Response**:

```json
{   "ticker": "AAPL",   "companyName": "Apple Inc.",   "currentPrice": 150.25,   "lastUpdated": "2024-12-11T10:00:00Z" }
```
---
#### All Stocks

**Endpoint**: `/all-stocks`  
**Method**: `GET`  
**Response**:
```json


{"AAPL": { "ticker": "AAPL", "currentPrice": 150.25 }, "TSLA": { "ticker": "TSLA", "currentPrice": 700.50 }...}
```

---

### Placing Orders

**Endpoint**: `/order`  
**Method**: `POST`  
**Request Body**:
```json
{   "quantity": 10,   "ticker": "AAPL",   "username": "example_user" }`

**Response**:
```
---
**Response:**
```json


{"message": "Successful buy order!",   "price": 150,   "ticker": "AAPL",   "quantity": 10}
```

**Note**: Ensure the account has sufficient balance for buy orders. For sell orders, simply set a negative quantity. 

---