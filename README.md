# GoAPI - A Beginner's Guide to Building REST APIs in Go

Welcome! This tutorial will walk you through a complete Go REST API project that retrieves user coin balances with authentication. Whether you're new to Go or building your first API, this project demonstrates essential concepts in a clean, organized way.

## 📚 What This Project Does

GoAPI is a simple but powerful REST API service that:
- Authenticates users with tokens
- Retrieves user coin balances
- Uses middleware for request validation
- Implements clean code architecture with separation of concerns

The API has one endpoint:
```
GET /account/coins?username=alex
Headers: Authorization: 123ABC
```

This endpoint checks if the user is authorized, then returns their coin balance.

---

## 🏗️ Project Structure

Let's understand how the code is organized:

```
goapi/
├── go.mod                          # Project dependencies
├── cmd/api/main.go                # Entry point - starts the server
├── api/api.go                     # Response/Request types & error handlers
├── internal/
│   ├── handlers/                  
│   │   ├── api.go                # Routes configuration
│   │   └── get_coin_balance.go   # Endpoint handler logic
│   ├── middleware/
│   │   └── authorization.go      # Authentication middleware
│   └── tools/
│       ├── database.go           # Database interface & setup
│       └── mockdb.go             # Mock database with test data
```

### **Directory Naming Convention**
- **`cmd/`** - Command line applications (entry points)
- **`api/`** - Public API types and handlers
- **`internal/`** - Private internal code (not exported to other packages)

This structure ensures clean boundaries between public and private code!

---

## 🔧 Dependencies

The project uses these Go packages:

```go
github.com/go-chi/chi v1.5.5           // Lightweight HTTP router
github.com/gorilla/schema v1.4.1       // Parse URL query parameters
github.com/sirupsen/logrus v1.9.3      // Structured logging
golang.org/x/sys                       // System utilities
```

---

## 🚀 How It Works - Step by Step

### **1. Starting the Server (`cmd/api/main.go`)**

```go
func main() {
    var r *chi.Mux = chi.NewRouter()
    handlers.Handler(r)
    http.ListenAndServe("localhost:8000", r)
}
```

**What happens:**
1. Create a new Chi router (a lightweight HTTP multiplexer)
2. Register all routes and middleware
3. Start listening on `localhost:8000`

Chi is like a traffic controller that directs HTTP requests to the right handlers.

---

### **2. Setting Up Routes (`internal/handlers/api.go`)**

```go
func Handler(r *chi.Mux) {
    r.Use(chimiddle.StripSlashes)  // Remove trailing slashes
    
    r.Route("/account", func(router chi.Router) {
        router.Use(middleware.Authorization)  // Protect this route
        router.Get("/coins", GetCoinBalance)   // Handle GET /account/coins
    })
}
```

**Breaking it down:**
- `r.Use()` - Apply global middleware (all requests)
- `r.Route()` - Define a route group with its own middleware
- `router.Use(middleware.Authorization)` - Only requests to `/account/*` need authentication
- `router.Get("/coins", GetCoinBalance)` - Map GET requests to our handler function

---

### **3. Authenticating Requests (`internal/middleware/authorization.go`)**

Middleware is code that runs **before** your handler. It's like a security checkpoint:

```go
func Authorization(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if username and token are provided
        var username string = r.URL.Query().Get("username")
        var token = r.Header.Get("Authorization")
        
        if username == "" || token == "" {
            api.RequestErrorHandler(w, UnAuthorizedError)
            return  // Stop here - don't proceed
        }
        
        // Get user's real token from database
        database, _ := tools.NewDatabase()
        loginDetails := (*database).GetUserLoginDetails(username)
        
        // Check if token matches
        if loginDetails == nil || (token != (*loginDetails).AuthToken) {
            api.RequestErrorHandler(w, UnAuthorizedError)
            return  // Stop here
        }
        
        // Token is valid! Call the next handler
        next.ServeHTTP(w, r)
    })
}
```

**Key concept:** If authentication fails, we return an error and never call `next.ServeHTTP()`. If it passes, we do!

---

### **4. Getting the Coin Balance (`internal/handlers/get_coin_balance.go`)**

This is the main handler that executes when all checks pass:

```go
func GetCoinBalance(w http.ResponseWriter, r *http.Request) {
    // Step 1: Parse the username from URL query
    var params = api.CoinBalanceParams{}
    var decoder *schema.Decoder = schema.NewDecoder()
    err := decoder.Decode(&params, r.URL.Query())
    
    // Step 2: Create a database connection
    database, err := tools.NewDatabase()
    
    // Step 3: Get the user's coins
    tokenDetails := (*database).GetUserCoins(params.Username)
    
    // Step 4: Send back the response
    var response = api.CoinBalanceResponse{
        Balance:    tokenDetails.Coins,
        StatusCode: http.StatusOK,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

**The flow:**
1. Extract the `username` parameter from the URL
2. Connect to the database
3. Fetch the user's coin balance
4. Return a JSON response with the data

---

### **5. Types and Responses (`api/api.go`)**

All the data structures are defined here:

```go
type CoinBalanceParams struct {
    Username string  // What we expect from the request
}

type CoinBalanceResponse struct {
    StatusCode int
    Balance    int64  // User's coins
}

type Error struct {
    StatusCode int
    Message    string
}
```

These types help Go validate the data and ensure type safety!

---

### **6. The Database (`internal/tools/database.go` and `internal/tools/mockdb.go`)**

We use an **interface** to separate the database logic from the rest of the code:

```go
type DatabaseInterface interface {
    GetUserLoginDetails(username string) *LoginDetails
    GetUserCoins(username string) *CoinDetails
    SetupDatabase() error
}
```

This is powerful because:
- We can swap the mock database with a real one without changing any other code
- It makes testing easier
- It follows the **Dependency Inversion Principle**

The mock database has test data:

```go
var mockLoginDetails = map[string]LoginDetails{
    "alex": {AuthToken: "123ABC", Username: "alex"},
    "maria": {AuthToken: "456DEF", Username: "maria"},
    "john": {AuthToken: "789GHI", Username: "john"},
}

var mockCoinDetails = map[string]CoinDetails{
    "alex": {Coins: 1000, Username: "alex"},
    "maria": {Coins: 2500, Username: "maria"},
    "john": {Coins: 500, Username: "john"},
}
```

---

## 📝 Example Request & Response

### **Request:**
```bash
curl -H "Authorization: 123ABC" "http://localhost:8000/account/coins?username=alex"
```

### **Successful Response (200 OK):**
```json
{
  "StatusCode": 200,
  "Balance": 1000
}
```

### **Failed Authentication (400 Bad Request):**
```bash
curl -H "Authorization: WRONG_TOKEN" "http://localhost:8000/account/coins?username=alex"
```

```json
{
  "StatusCode": 400,
  "Message": "Invalid username or token."
}
```

---

## 🎯 Key Concepts Explained

### **1. Middleware**
Middleware is a function that wraps your handler. It runs **before** the handler and can:
- Validate requests (authentication, validation)
- Modify requests
- Stop processing if needed

Think of it like a filter on a coffee maker - some things pass through, some don't.

### **2. Interfaces**
An interface defines a contract - what methods an object must have:

```go
type DatabaseInterface interface {
    GetUserLoginDetails(username string) *LoginDetails
    GetUserCoins(username string) *CoinDetails
    SetupDatabase() error
}
```

Any type that implements these three methods satisfies the interface. This allows us to:
- Switch databases without changing handler code
- Mock databases for testing
- Keep code loosely coupled

### **3. Pointers**
Notice `*DatabaseInterface` and `*LoginDetails`. The asterisk means "pointer" - it's a reference to the actual data:

```go
database, err := tools.NewDatabase()  // Returns a pointer
loginDetails := (*database).GetUserLoginDetails(username)  // Dereference with *
```

Pointers are efficient for large data structures!

### **4. Error Handling**
Go handles errors explicitly:

```go
if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
}
```

This is better than throwing exceptions because you always know where errors can happen.

---

## 🚀 Running the Project

### **1. Clone and navigate to the project:**
```bash
cd d:\goapi
```

### **2. Install dependencies:**
```bash
go mod download
```

### **3. Run the server:**
```bash
go run cmd/api/main.go
```

You should see:
```
Starting GO API service....
  _________    ___   ___  ____
 / ___/ __ \  / _ | / _ \/  _/
/ (_ / /_/ / / __ |/ ___// /  
\___/\____/ /_/ |_/_/  /___/
```

### **4. Test the API:**
```bash
# Valid request
curl -H "Authorization: 123ABC" "http://localhost:8000/account/coins?username=alex"

# Invalid token
curl -H "Authorization: WRONG" "http://localhost:8000/account/coins?username=alex"

# Missing credentials
curl "http://localhost:8000/account/coins?username=alex"
```

---

## 🎓 Learning Path

If you're new to Go, study in this order:

1. **Read `cmd/api/main.go`** - Understand how the server starts
2. **Read `internal/handlers/api.go`** - Learn how routes are configured
3. **Read `internal/middleware/authorization.go`** - See how middleware works
4. **Read `internal/handlers/get_coin_balance.go`** - Understand request handling
5. **Read `internal/tools/database.go`** - Learn about interfaces
6. **Read `api/api.go`** - See how types are used for validation

---

## 📚 Next Steps

Now that you understand this API, you can:

- **Add more endpoints** - Create new handler functions
- **Connect a real database** - Replace mockDB with PostgreSQL, MongoDB, etc.
- **Add validation** - Check that coins is a positive number
- **Add logging** - Track important events
- **Write tests** - Create unit tests for handlers
- **Add caching** - Improve performance with Redis
- **Deploy** - Put it on a server for the world to use

---

## 💡 Best Practices Used

✅ **Clean Architecture** - Separation of concerns (handlers, middleware, tools)  
✅ **Interfaces** - Database abstraction makes code flexible  
✅ **Error Handling** - Explicit error checks throughout  
✅ **Middleware Pattern** - Reusable request processing  
✅ **Type Safety** - Strong typing with structs  
✅ **Logging** - Track errors and important events  

---

## 📖 Resources

- [Go Documentation](https://golang.org/doc/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [HTTP Package Guide](https://golang.org/pkg/net/http/)

Happy coding! 🎉
