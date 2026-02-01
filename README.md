# Kasir API

A RESTful API for managing products and categories in a cashier system, built with Go standard library using a clean
architecture pattern.

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

- ✅ RESTful API dengan metode HTTP standar
- ✅ Clean architecture (entity, data, service, handler, router)
- ✅ In-memory storage yang bisa diganti (swappable)
- ✅ Validasi input (name wajib, price > 0, stock >= 0)
- ✅ JSON response terstandardisasi
- ✅ Error handling dengan custom error types
- ✅ Thread-safe operations dengan mutex
- ✅ Hanya menggunakan Go standard library (`net/http`)
- ✅ Integer IDs (auto-increment) dan integer prices

## API Endpoints

### Health Check

```
GET /health
```

Mengembalikan status kesehatan API.

Response (sesuai implementasi helper/response.go):

```json
{
  "status": "OK",
  "message": "API Running"
}
```

### Product Endpoints

#### Get All Products

```
GET /api/products
```

Mengambil semua produk.

Response (contoh):

```json
{
  "status": "OK",
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
GET /api/productss/{id}
```

Mengambil produk berdasarkan ID.

Response:

```json
{
  "status": "OK",
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
POST /api/products
Content-Type: application/json
```

Request Body:

```json
{
  "name": "Laptop",
  "price": 15000000,
  "stock": 10
}
```

Response:

```json
{
  "status": "OK",
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
PUT /api/products/{id}
Content-Type: application/json
```

Request Body:

```json
{
  "name": "Laptop Gaming",
  "price": 20000000,
  "stock": 5
}
```

Response:

```json
{
  "status": "OK",
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
DELETE /api/products/{id}
```

Response:

```json
{
  "status": "OK",
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

Request Body:

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

Request Body:

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

Semua error mengikuti format standar dari helper/response.go:

```json
{
  "status": "ERROR",
  "message": "Failed to create product"
}
```

Contoh pesan validasi yang mungkin muncul:

- `Invalid JSON` — payload tidak dapat di-decode
- `Invalid Product ID` atau `Invalid Category ID` — ID tidak valid
- `Product not found` atau `Category not found` — resource tidak ditemukan
- `Failed to create/update/delete product/category` — kegagalan bisnis/validasi

### HTTP Status Codes

- `200 OK` — Berhasil untuk GET, PUT, DELETE
- `201 Created` — Berhasil untuk POST
- `400 Bad Request` — Validasi error atau input tidak valid
- `404 Not Found` — Resource tidak ditemukan
- `405 Method Not Allowed` — HTTP method tidak didukung
- `500 Internal Server Error` — Error server

## Installation & Running

### Prasyarat

- Go 1.25.6 atau lebih baru

### Build

```bash
go build -o kasir-api .
```

### Run

```bash
./kasir-api
```

Atau menggunakan `go run`:

```bash
go run .
```

Server akan berjalan di port 8080 secara default. Ubah port dengan environment variable `PORT`:

```bash
PORT=3000 ./kasir-api
```

Catatan: Storage in-memory akan di-reset setiap kali server restart; ID auto-increment dimulai dari 1 pada sesi baru.

## Testing

Skrip `test_api.sh` tersedia untuk menguji seluruh endpoint API secara end-to-end.

### Prasyarat

- `curl` (wajib)
- `jq` (opsional, untuk menampilkan JSON dengan rapi)

### Cara Menjalankan

1. Jalankan server terlebih dahulu (lihat bagian Run). Pastikan base URL sesuai.
2. Jalankan skrip:

```bash
./test_api.sh
```

Secara default, skrip akan menggunakan `http://localhost:8080`. Anda bisa mengubah Base URL dengan argumen pertama:

```bash
./test_api.sh http://localhost:3000
```

### Yang Diuji

Skrip akan menjalankan urutan berikut:

- Health check: `GET /health`
- Category: create beberapa kategori, get all, get by ID, update, dan validasi error (invalid ID/JSON, not found)
- Product: create beberapa produk, get all, get by ID, update, dan validasi error (invalid ID/JSON, not found)
- Delete: hapus produk dan kategori tertentu, lalu verifikasi 404 pada akses selanjutnya

Skrip menilai berdasarkan HTTP status code (2xx sukses, 4xx untuk error yang diharapkan, atau exact codes seperti 404).
Di akhir, skrip menampilkan total, passed/failed, dan success rate, serta exit code 0 jika semua lulus.

### Troubleshooting

- Pastikan server berjalan dan dapat diakses di Base URL.
- Jika menggunakan port custom, sesuaikan argumen `./test_api.sh <BASE_URL>`.
- Bersihkan state dengan me-restart server untuk in-memory storage yang fresh.
- Jika `jq` tidak terpasang, output body akan ditampilkan apa adanya.

## Swappable Storage

Layer storage menggunakan pola interface, sehingga mudah diganti implementasinya:

```go
type ProductStorage interface {
GetAll() ([]*entity.Product, error)
GetByID(id int) (*entity.Product, error)
Create(product *entity.Product) error
Update(product *entity.Product) error
Delete(id int) error
}
```

Implementasi saat ini: `MemoryStorage` (in-memory dengan mutex untuk thread safety)

Untuk menambah implementasi storage baru (mis. database):

1. Buat file baru di `data/` (mis. `postgres_storage.go`)
2. Implementasikan interface storage
3. Perbarui `main()` untuk menggunakan storage baru:

```go
storage := data.NewPostgresStorage(connectionString)
```

## Project Structure Details

### Entity Layer

- `product.go`: domain model Product dengan JSON tags
- `category.go`: domain model Category dengan JSON tags
- `errors.go`: custom error types untuk business logic

### Data Layer

- `product_storage.go`: interface storage untuk produk
- `category_storage.go`: interface storage untuk kategori
- `category_memory_storage.go`: adapter kategori di atas MemoryStorage
- `memory_storage.go`: implementasi in-memory dengan operasi thread-safe

### Service Layer

- `product_service.go`: business logic untuk produk (validasi, CRUD)
- `category_service.go`: business logic untuk kategori (validasi, CRUD)

### Handler Layer

- `product_handler.go`: HTTP handler untuk produk
- `category_handler.go`: HTTP handler untuk kategori

### Helper Layer

- `response.go`: format JSON response standar (`status`, `message`, `data`)

### Router Layer

- `router.go`: routing HTTP dengan dispatch berdasarkan method

## License

GPL
