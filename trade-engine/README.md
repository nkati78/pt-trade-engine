# Trade Engine
The Trade engine is responsible for users, orders, and trades. It is the core of the exchange. It is responsible for matching orders and executing trades. It is also responsible for maintaining the order book and the user's balances.

## Installation
```bash
brew install docker
brew install docker-compose
```
## Building the Trade Engine
```bash
make build
```

## Running the Trade Engine
```bash
make run
```

Ports are exposed on `localhost:8080`

## Routes
### Orders
- `GET /orders` - Get all orders for a user, Authorization header with a valid JWT token is required.
```javascript
    // curl -x GET http://localhost:8080/orders -H "Authorization: $JWT_TOKEN"
```
- `POST /orders` - Create a new order, Authorization header with a valid JWT token is required.
```json 
 // curl -x POST http://localhost:8080/orders -H "Authorization: $JWT_TOKEN"
  {
    "symbol": "AAPL",
    "side": "BUY",
    "price": 10000,
    "quantity": 1,
    "type": "MARKET"
  }
 ```

### Users
- `POST /register` - Create a new user, returns a JWT token
```json
  // curl -x POST http://localhost:8080/register 
  {
	"username": "jubjub",
	"firstName": "josh", 
	"lastName": "horecny",
	"email": "josh@paperthesis.com",
	"password": "password!"
  }
```
- `POST /login` - Login a user, returns a JWT token
```json
  // curl -x POST http://localhost:8080/login
  {   
    "username": "josh@paperthesis.com",
    "password": "password!"
  }
```
