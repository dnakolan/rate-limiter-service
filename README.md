# rate-limiter-service
A Go service that enforces per user rate limits

# Features
* Two Rate Limiting Algorithms:
** Fixed Window
** Sliding Window
* HTTP Middleware - Apply rate limiting to any endpoint
* Configuration API - Manage rate limits via REST
* In-Memory Storage - Keep it simple, no external dependencies
* Basic Cleanup - Remove expired data periodically

# Tech Stack
* Language: Go 1.24+
* Router: gin
* UUIDs: github.com/google/uuid
* Testing: Go standard library

# Project Structure
```
rate-limiter-service/
├── models/
├──── cmd/
│     └── main.go     // application entry point
├── handlers/
│   └── limits.go     // http handlers for /limits
│   └── health.go     // http handlers for /health
├── models/
│   └── limits.go     // data models and validation for /limits
├── services/
│   └── limits.go     // supporting service for /limits endpoints
├── storage/
│   └── memory.go     // in memory store of event data
└── build.sh          // builds docker images for the application
└── Dockerfile        // core application dependencies
└── Dockerfile.deps   // isolated base image to speed up docker build
```

# RateLimit Schema
```
{
  "id": "api-users",
  "limit": 100,
  "window": "1m",
  "algorithm": "sliding_window",
  "key_pattern": "user:{{user_id}}"
}
```

# Running the Service
First run the included build.sh script to build the container images
```
./build.sh
```

Then start the application in docker with the following command
```
docker run -d -p 8080:8080 --name rate-service rate-limiting-service
```

# Example Usage (cURL)
## POST /limits - Create rate limit rule
```
curl -X POST http://localhost:8080/limits \
  -H "Content-Type: application/json" \
  -d '{
        "id": "api-users",
        "limit": 100,
        "window": "1m",
        "algorithm": "sliding_window",
        "key_pattern": "user:{{user_id}}"
      }'
```

## GET /limits - List all rules 
```
curl http://localhost:8080/limits

[
  {
    "id": "api-users",
    "limit": 100,
    "window": "1m",
    "algorithm": "sliding_window",
    "key_pattern": "user:{{user_id}}"
  }
]
```

## PUT /limits/:id - Update rule
```
curl -X PUT http://localhost:8080/limits \
  -H "Content-Type: application/json" \
  -d '{
        "id": "api-users",
        "limit": 200,
        "window": "1m",
        "algorithm": "sliding_window",
        "key_pattern": "user:{{user_id}}"
      }'
```

## DELETE /limits/:id - Delete rule
```
curl -X PUT http://localhost:8080/limits/api-users
```

## GET /check/:userID - Check if user can make request
```
curl http://localhost:8080/check/api-users
```

## POST /middleware - Apply rate limiting to requests
```

# Design Considerations
* Dependency Injection is used for loose coupling between components.
* Interface-Driven Architecture enables testability and future extensibility (e.g., database-backed repo).
* Validation is handled at the request model level to separate concerns cleanly.
* The service layer enforces any domain-specific business rules.
* Websocket interface for processing batches of events

# Tests
`go test ./...`
Tests cover handler logic, service behavior, and in-memory repo operations.

# Future Improvements / Next Steps
TBD

# Time Spent
TBD

# Author
David Nakolan - david.nakolan@gmail.com
