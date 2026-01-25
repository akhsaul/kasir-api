# Kasir API

A RESTful API for managing products and categories in a cashier system, built with Go standard library using a clean architecture pattern.

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

```
kasir-api/
├── entity/         # Domain entities and errors
│   ├── product.go
│   ├── category.go
│   └── errors.go
├── data/           # Data access layer (storage interface + implementations)
│   ├── product_storage.go
│   ├── category_storage.go
│   ├── category_memory_storage.go
│   └── memory_storage.go
├── service/        # Business logic layer
│   ├── product_service.go
│   └── category_service.go
├── handler/        # HTTP request handlers
│   ├── product_handler.go
│   └── category_handler.go
├── helper/         # Utility functions
│   └── response.go
├── router/         # HTTP routing
│   └── router.go
└── kasir_api.go    # Application entry point
```

## Features

- ✅ RESTful API with standard HTTP methods
- ✅ Clean architecture (entity, data, service, handler, router layers)
- ✅ In-memory storage with swappable interface
- ✅ Input validation (name required, price > 0, stock >= 0)
- ✅ Standardized JSON responses
- ✅ Error handling with custom error types
- ✅ Thread-safe operations with mutex
- ✅ Uses only Go standard library (`net/http`)
- ✅ Integer IDs (auto-increment) and integer prices

## API Endpoints

### Health Check
```
GET /health
```
Returns API health status.

**Response:**
```json
{
  "status": "healthy"
}
```

### Product Endpoints

#### Get All Products
```
GET /api/product
```
Retrieves all products.

**Response (example):**
```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Laptop",
      "price": 15000000,
      "stock": 10
    }
  ]
}
```

#### Get Product by ID
```
GET /api/product/{id}
```
Retrieves a specific product by ID.

**Response:**
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "name": "Laptop",
    "price": 15000000,
    "stock": 10
  }
}
```

#### Create Product
```
POST /api/product
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Laptop",
  "price": 15000000,
  "stock": 10
}
```

**Response:**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "name": "Laptop",
    "price": 15000000,
    "stock": 10
  }
}
```

#### Update Product
```
PUT /api/product/{id}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Laptop Gaming",
  "price": 20000000,
  "stock": 5
}
```

**Response:**
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "name": "Laptop Gaming",
    "price": 20000000,
    "stock": 5
  }
}
```

#### Delete Product
```
DELETE /api/product/{id}
```

**Response:**
```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

### Category Endpoints

#### Get All Categories
```
GET /api/categories
```

#### Get Category by ID
```
GET /api/categories/{id}
```

#### Create Category
```
POST /api/categories
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

#### Update Category
```
PUT /api/categories/{id}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Electronics & Gadgets",
  "description": "All electronic devices, gadgets and accessories"
}
```

#### Delete Category
```
DELETE /api/categories/{id}
```

## Error Responses

All errors follow a standard format:

```json
{
  "success": false,
  "message": "Failed to create product",
  "error": "name is required"
}
```

### Validation Errors

- `name is required` - Name cannot be empty
- `price must be greater than 0` - Price must be positive (integer)
- `stock must be greater than or equal to 0` - Stock cannot be negative
- `id is required` - ID parameter is missing or invalid

### HTTP Status Codes

- `200 OK` - Successful GET, PUT, DELETE
- `201 Created` - Successful POST
- `400 Bad Request` - Validation error or invalid input
- `404 Not Found` - Resource not found
- `405 Method Not Allowed` - HTTP method not supported
- `500 Internal Server Error` - Server error

## Installation & Running

### Prerequisites
- Go 1.25.6 or later

### Build
```bash
go build -o kasir-api .
```

### Run
```bash
./kasir-api
```

Or using `go run`:
```bash
go run .
```

The server will start on port 8080 by default. You can change the port using the `PORT` environment variable:
```bash
PORT=3000 ./kasir-api
```

## Swappable Storage

The storage layer uses an interface pattern, making it easy to swap implementations:

```go
type ProductStorage interface {
    GetAll() ([]*entity.Product, error)
    GetByID(id int) (*entity.Product, error)
    Create(product *entity.Product) error
    Update(product *entity.Product) error
    Delete(id int) error
}
```

Current implementation: `MemoryStorage` (in-memory with mutex for thread safety)

To add a new storage implementation (e.g., database):
1. Create a new file in `data/` (e.g., `postgres_storage.go`)
2. Implement the storage interface
3. Update `main()` to use the new storage:
```go
storage := data.NewPostgresStorage(connectionString)
```

## Project Structure Details

### Entity Layer
- **product.go**: Product domain model with JSON tags
- **category.go**: Category domain model with JSON tags
- **errors.go**: Custom error types for business logic

### Data Layer
- **product_storage.go**: Product storage interface
- **category_storage.go**: Category storage interface
- **category_memory_storage.go**: Category adapter on MemoryStorage
- **memory_storage.go**: In-memory implementation with thread-safe operations

### Service Layer
- **product_service.go**: Business logic for products (validation, CRUD)
- **category_service.go**: Business logic for categories (validation, CRUD)

### Handler Layer
- **product_handler.go**: HTTP request/response handling for products
- **category_handler.go**: HTTP request/response handling for categories

### Helper Layer
- **response.go**: Standardized JSON response formatting

### Router Layer
- **router.go**: HTTP routing with method-based dispatching

## License

GPL
