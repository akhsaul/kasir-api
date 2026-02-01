Tutorial Lengkap - Sesi 1: Build Kasir API dari Nol

## \#🎯 Apa yang Akan Kita Bangun?

Hari ini kita akan bikin **Kasir API** \- sistem kasir sederhana yang bisa:

- Kelola data produk (CRUD lengkap)

- Terima request lewat HTTP

- Kirim response dalam format JSON

- Siap di-deploy ke cloud (gratis!)

Yang keren dari Go? Binary-nya **kecil banget** (kurang dari 15MB) dan **gak butuh runtime** kayak PHP atau Node.js.
Langsung jalan!

* * *

## \#🚀 Persiapan Awal

### \#Yang Harus Kamu Punya

1. **Go terinstall** (cek: `go version`)

2. **Text editor** (VSCode recommended)

3. **Terminal/Command Prompt**

4. **Git** (untuk deploy nanti)

### \#Buat Project Baru

```
# Bikin folder project
mkdir kasir-api
cd kasir-api

# Initialize Go module
go mod init kasir-api

# Bikin file main.go
touch main.go  # Mac/Linux
# atau
type nul > main.go  # Windows
```

**Kenapa**`go mod init` **?** Ini kayak `npm init` atau `composer init` \- bikin file konfigurasi untuk manage
dependencies. File `go.mod` akan dibuat otomatis.

* * *

## \#📦 CODE-01: Package & Import

Buka `main.go`, ketik ini:

```
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)
```

### \#Penjelasan Detail

`package main` \- Ini "pintu masuk" aplikasi Go.

- Setiap file Go harus punya `package`

- `package main` = executable program (bisa dijalankan)

- Package lain = library (dipanggil dari tempat lain)

**Import** \- Library yang kita butuhkan:

- `encoding/json` \- Encode/decode JSON (buat API response)

- `fmt` \- Print ke console (`fmt.Println`)

- `net/http` \- HTTP server & handling

- `strconv` \- Convert string ke number (untuk ID dari URL)

- `strings` \- Manipulasi string (trim, split, dll)

* * *

## \#🏗️ CODE-02: Produk Struct

Tambahkan setelah import:

```
// Produk represents a product in the cashier system
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}
```

### \#Penjelasan Detail

**Struct** = Blueprint data, mirip class tapi lebih simple.

**Breakdown field:**

- `ID int` \- ID produk (integer)

- `Nama string` \- Nama produk (text)

- `Harga int` \- Harga produk (integer, dalam rupiah)

- `Stok int` \- Jumlah stok

**Tag**`json:"..."` \- Ini magic! Ketika struct di-convert ke JSON, field name akan jadi lowercase:

```
{
  "id": 1,        // bukan "ID"
  "nama": "...",  // bukan "Nama"
  "harga": 3500,  // bukan "Harga"
  "stok": 100     // bukan "Stok"
}
```

Tanpa tag, output JSON-nya bakal `{"ID": 1, "Nama": "..."}` \- gak bagus untuk API.

**Kenapa**`Harga int` **?** Kita pakai integer (bukan float) karena rupiah gak punya desimal. Lebih simple & gak ada
masalah floating point!

* * *

## \#💾 CODE-03: In-Memory Storage

```
// In-memory storage (sementara, nanti ganti database)
var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
}
```

### \#Penjelasan Detail

`var produk []Produk` \- Variabel global berisi array of Produk.

**Kenapa pakai array?** Untuk sesi 1, kita fokus ke HTTP handling dulu. Sesi 2 nanti baru pakai **SQLite database**.

**Data dummy** \- Produk warung Indo klasik:

- Indomie Godog (favorit!)

- Vit (air mineral lokal)

- Kecap (buat masak)

**Catatan Penting:** Data ini **hilang** kalau server restart. Makanya namanya "in-memory". Database persisten nanti di
Sesi 2.

* * *

## \#🏥 CODE-04: Health Check Endpoint

Ini endpoint pertama kita - buat ngecek server hidup atau nggak:

```
// localhost:8080/health
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
})
```

### \#Penjelasan Detail

**Function signature:**

- `http.HandleFunc` \- Daftarkan route `/health`

- `func(w http.ResponseWriter, r *http.Request)` \- Anonymous function (function tanpa nama)

- `w` \- Tempat nulis response ke user

- `r` \- Data request dari user (method, headers, body, dll)

**Step by step:**

1. **Set header:**

```
w.Header().Set("Content-Type", "application/json")
```

Kasih tau browser: "Ini JSON ya, bukan HTML"

2. **Bikin & kirim response:**

```
json.NewEncoder(w).Encode(map[string]string{
       "status":  "OK",
       "message": "API Running",
})
```

`map[string]string` = object/dictionary dengan key & value string, langsung di-encode jadi JSON

**Output:**

```
{
  "status": "OK",
  "message": "API Running"
}
```

* * *

## \#📋 CODE-05: Get All Produk & Create Produk

```
// GET localhost:8080/api/produk
// POST localhost:8080/api/produk
http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produk)

	} else if r.Method == "POST" {
		// baca data dari request
		var produkBaru Produk
		err := json.NewDecoder(r.Body).Decode(&produkBaru)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// masukkin data ke dalam variable produk
		produkBaru.ID = len(produk) + 1
		produk = append(produk, produkBaru)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(produkBaru)
	}
})
```

### \#Penjelasan Detail

**Multi-method handler** \- 1 function handle GET & POST!

**Checking method dengan if-else:**

**1\. GET** \- Return semua produk:

```
if r.Method == "GET" {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(produk)
}
```

Array `produk` langsung jadi JSON array

**2\. POST** \- Create produk baru:

**Step 1 - Decode JSON dari body:**

```
var produkBaru Produk
err := json.NewDecoder(r.Body).Decode(&produkBaru)
```

- Baca JSON dari request body

- Convert ke struct Produk

- `&produkBaru` \- pointer, biar langsung isi data ke variable

**Step 2 - Error handling:**

```
if err != nil {
    http.Error(w, "Invalid request", http.StatusBadRequest)
    return
}
```

Kalau JSON rusak → return 400

**Step 3 - Auto-increment ID & append:**

```
produkBaru.ID = len(produk) + 1
produk = append(produk, produkBaru)
```

- Generate ID otomatis

- Tambah ke array

**Step 4 - Response:**

```
w.WriteHeader(http.StatusCreated) // 201
json.NewEncoder(w).Encode(produkBaru)
```

Return produk yang baru dibuat

**Test GET:**

```
curl http://localhost:8080/api/produk
```

**Response:**

```
[\
  {\
    "id": 1,\
    "nama": "Indomie Godog",\
    "harga": 3500,\
    "stok": 10\
  },\
  {\
    "id": 2,\
    "nama": "Vit 1000ml",\
    "harga": 3000,\
    "stok": 40\
  },\
  {\
    "id": 3,\
    "nama": "kecap",\
    "harga": 12000,\
    "stok": 20\
  }\
]
```

**Test POST:**

```
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Kopi Kapal Api",
    "harga": 2500,
    "stok": 200
  }'
```

**Response:**

```
{
  "id": 4,
  "nama": "Kopi Kapal Api",
  "harga": 2500,
  "stok": 200
}
```

* * *

## \#🔍 CODE-06: Get Produk by ID

Ini yang agak tricky - ambil ID dari URL path!

```
func getProdukByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/produk/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// Cari produk dengan ID tersebut
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}
```

### \#Penjelasan Detail

**URL Path Parsing** \- Ini yang paling tricky:

1. **URL:**`/api/produk/123`

2. **TrimPrefix:** Hilangkan `/api/produk/` → dapat `"123"`

3. **Atoi:** Convert `"123"` string → `123` integer

```
idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
// URL: /api/produk/123 -> idStr = "123"
// URL: /api/produk/abc -> idStr = "abc"

id, err := strconv.Atoi(idStr)
// "123" -> 123 ✅
// "abc" -> error ❌
```

**Loop & Search:**

```
for _, p := range produk {
    if p.ID == id {
        // Found! Return produk
        json.NewEncoder(w).Encode(p)
        return
    }
}
```

- `range produk` \- Loop semua produk

- `_` \- Ignore index (gak butuh)

- `p` \- Current produk di loop

- Return langsung kalau ketemu

**Not Found:**

```
http.Error(w, "Produk belum ada", http.StatusNotFound)
```

Kalau sampai sini = loop selesai, gak ketemu → 404

**Test:**

```
# Found
curl http://localhost:8080/api/produk/1

# Not found
curl http://localhost:8080/api/produk/999
```

**Response (Found):**

```
{
  "id": 1,
  "nama": "Indomie Godog",
  "harga": 3500,
  "stok": 10
}
```

**Response (Not Found):**

```
404 Produk belum ada
```

* * *

## \#✏️ CODE-07: Update Produk

```
// PUT localhost:8080/api/produk/{id}
func updateProduk(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop produk, cari id, ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}
```

### \#Penjelasan Detail

**Combine 2 Logic:**

1. Parse ID dari URL (kayak GET by ID)

2. Decode body JSON (kayak CREATE)

**Step by step:**

**1\. Parse ID dari URL:**

```
idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
id, err := strconv.Atoi(idStr)
```

**2\. Decode update data:**

```
var updateProduk Produk
err = json.NewDecoder(r.Body).Decode(&updateProduk)
```

**3\. Loop & Update:**

```
for i := range produk {
    if produk[i].ID == id {
        updateProduk.ID = id      // Keep ID!
        produk[i] = updateProduk  // Update di index i
        return
    }
}
```

**Kenapa**`for i := range produk` **?**

- `i` = index di array (0, 1, 2, ...)

- `produk[i] = updateProduk` \- Replace produk di index tersebut

- **PENTING:** Keep ID tetap sama! User gak bisa ganti ID

**Test:**

```
curl -X PUT http://localhost:8080/api/produk/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Indomie Goreng Jumbo",
    "harga": 4000,
    "stok": 150
  }'
```

**Response:**

```
{
  "id": 1,  // ID tetap sama!
  "nama": "Indomie Goreng Jumbo",
  "harga": 4000,
  "stok": 150
}
```

* * *

## \#🗑️ CODE-08: Delete Produk

```
func deleteProduk(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// loop produk cari ID, dapet index yang mau dihapus
	for i, p := range produk {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}
```

### \#Penjelasan Detail

**Delete Trick** \- Ini yang agak magic:

```
produk = append(produk[:i], produk[i+1:]...)
```

**Visualisasi:**

```
Array sebelum delete:
Index: 0     1     2     3
Value: [A] - [B] - [C] - [D]

Mau delete index 1 (B):
produk[:1]  = [A]          // Sebelum index 1
produk[2:]  = [C, D]       // Sesudah index 1
Append = [A] + [C, D] = [A, C, D]

Array sesudah delete:
Index: 0     1     2
Value: [A] - [C] - [D]
```

**Breakdown:**

- `produk[:i]` \- Semua element **sebelum** index i

- `produk[i+1:]` \- Semua element **sesudah** index i

- `append(...)` \- Gabungkan keduanya

- `...` \- Spread operator (buka array jadi individual elements)

**Response:**

```
json.NewEncoder(w).Encode(map[string]string{
    "message": "sukses delete",
})
```

Beda dari CRUD lain, DELETE gak return data. Cukup success message.

**Test:**

```
curl -X DELETE http://localhost:8080/api/produk/1
```

**Response:**

```
{
  "message": "sukses delete"
}
```

* * *

## \#🚦 CODE-09: Main Function & Routing

Sekarang hubungkan semua handler ke routes:

```
func main() {
	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			// baca data dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
```

### \#Penjelasan Detail

**Routing dengan Anonymous Function:** Semua handler pakai `func(w http.ResponseWriter, r *http.Request)` langsung - gak
perlu bikin function terpisah untuk simple logic.

**PENTING - Slash di akhir!**

```
http.HandleFunc("/api/produk", ...)   // tanpa slash
http.HandleFunc("/api/produk/", ...)  // dengan slash
```

- `/api/produk` (tanpa slash) = **exact match**

    - Match: `/api/produk` ✅

    - Tidak match: `/api/produk/1` ❌
- `/api/produk/` (dengan slash) = **prefix match**

    - Match: `/api/produk/` ✅

    - Match: `/api/produk/1` ✅

    - Match: `/api/produk/123` ✅

Makanya kita butuh 2 route handler berbeda!

**Handler Functions:** Untuk GET/PUT/DELETE by ID, kita panggil function terpisah (getProdukByID, updateProduk,
deleteProduk) karena logic-nya lebih kompleks.

**Start Server:**

```
http.ListenAndServe(":8080", nil)
```

- `:8080` \- Port (bisa ganti ke :3000, :5000, dll)

- `nil` \- Default HTTP handler (gak pakai framework)

**Run:**

```
go run main.go
```

**Output:**

```
Server running di localhost:8080
```

Server running! Buka browser: `http://localhost:8080/health`

### \#Available Endpoints

```
GET    /health              - Health check
GET    /api/produk          - Get all produk
POST   /api/produk          - Create produk
GET    /api/produk/{id}     - Get produk by ID
PUT    /api/produk/{id}     - Update produk
DELETE /api/produk/{id}     - Delete produk
```

* * *

## \#🧪 Testing Endpoints

### \#Pakai cURL (Command Line)

**Health check:**

```
curl http://localhost:8080/health
```

**Get all produk:**

```
curl http://localhost:8080/api/produk
```

**Create produk:**

```
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Kopi Kapal Api",
    "harga": 2500,
    "stok": 200
  }'
```

**Get produk by ID:**

```
curl http://localhost:8080/api/produk/1
```

**Update produk:**

```
curl -X PUT http://localhost:8080/api/produk/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Indomie Goreng Jumbo",
    "harga": 4000,
    "stok": 150
  }'
```

**Delete produk:**

```
curl -X DELETE http://localhost:8080/api/produk/1
```

### \#Pakai GUI Tool

Lebih enak pakai GUI? Install salah satu:

- **Postman** \- Popular, feature lengkap

- **Thunder Client** \- VSCode extension, ringan

- **Bruno** \- Open source, modern

Tutorial install ada di dokumen terpisah!

* * *

## \#🏗️ Build Binary

Go bisa di-compile jadi **single binary** \- gak perlu runtime!

### \#Build Standar

```
go build -o kasir-api
```

- Output: `kasir-api` (Mac/Linux) atau `kasir-api.exe` (Windows)

- Size: ~10-15MB

- No dependencies!

### \#Build Production (Smaller)

```
go build -ldflags="-s -w" -o kasir-api
```

- `-s` \- Strip symbol table

- `-w` \- Strip debug info

- Size: ~7-10MB (lebih kecil!)

### \#Cross-Compilation

**Build untuk Windows (dari Mac/Linux):**

```
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api.exe
```

**Build untuk Linux (dari Windows/Mac):**

```
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api
```

**Build untuk Mac (dari Windows/Linux):**

```
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api
```

Keren kan? 1 command, bisa build untuk semua OS!

### \#Jalankan Binary

```
# Mac/Linux
./kasir-api

# Windows
kasir-api.exe
```

Server langsung jalan - gak perlu `go run` lagi!

* * *

## \#☁️ Deploy ke Railway (Gratis!)

### \#Persiapan

**1\. Bikin**`.gitignore` **:**

```
kasir-api
kasir-api.exe
*.log
```

**2\. Push ke GitHub:**

```
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/username/kasir-api.git
git push -u origin main
```

### \#Deploy di Railway

1. Buka [railway.app](https://railway.app/)

2. Login dengan GitHub

3. Click **"New Project"**

4. Pilih **"Deploy from GitHub repo"**

5. Select repo `kasir-api`

6. Railway auto-detect Go → auto-deploy!

7. Dapat URL: `https://kasir-api-production.up.railway.app`

**Test production:**

```
curl https://kasir-api-production.up.railway.app/health
```

Works! 🎉

**Environment Variables:**

- Railway auto-set `PORT` variable

- Bisa tambah variable di dashboard

**Alternative ke zeabur baca dokumen [ini](https://docs.kodingworks.io/s/8b1de4f2-4642-4ac5-b068-7c9db2764c0e)**

* * *

## \#🎓 Apa yang Sudah Kita Pelajari?

✅ Go basics (package, import, struct, function)

✅ HTTP handling (request, response, routing)

✅ JSON encoding/decoding

✅ CRUD operations (in-memory)

✅ URL path parsing

✅ Error handling

✅ Build & deployment

* * *

## \#💡 Tips & Best Practices

**1\. Error Handling:**

```
if err != nil {
    // ALWAYS handle error!
    http.Error(w, "Error message", http.StatusBadRequest)
    return
}
```

**2\. HTTP Status Codes:**

- 200 OK - Success

- 201 Created - Resource created

- 400 Bad Request - Invalid input

- 404 Not Found - Resource not found

- 405 Method Not Allowed - Wrong HTTP method

- 500 Internal Server Error - Server error

**3\. Response Format:** Konsisten pakai JSON:

```
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```

* * *

0 / 0

### Contents

01. [🎯 Apa yang Akan Kita Bangun?](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%8E%AF-apa-yang-akan-kita-bangun)
02. [🚀 Persiapan Awal](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%9A%80-persiapan-awal)
03. [Yang Harus Kamu Punya](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-yang-harus-kamu-punya)
04. [Buat Project Baru](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-buat-project-baru)
05. [📦 CODE-01: Package & Import](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%93%A6-code-01-package-and-import)
06. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail)
07. [🏗️ CODE-02: Produk Struct](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%8F%97%EF%B8%8F-code-02-produk-struct)
08. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-1)
09. [💾 CODE-03: In-Memory Storage](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%92%BE-code-03-in-memory-storage)
10. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-2)
11. [🏥 CODE-04: Health Check Endpoint](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%8F%A5-code-04-health-check-endpoint)
12. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-3)
13. [📋 CODE-05: Get All Produk & Create Produk](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%93%8B-code-05-get-all-produk-and-create-produk)
14. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-4)
15. [🔍 CODE-06: Get Produk by ID](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%94%8D-code-06-get-produk-by-id)
16. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-5)
17. [✏️ CODE-07: Update Produk](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%E2%9C%8F%EF%B8%8F-code-07-update-produk)
18. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-6)
19. [🗑️ CODE-08: Delete Produk](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%97%91%EF%B8%8F-code-08-delete-produk)
20. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-7)
21. [🚦 CODE-09: Main Function & Routing](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%9A%A6-code-09-main-function-and-routing)
22. [Penjelasan Detail](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-penjelasan-detail-8)
23. [Available Endpoints](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-available-endpoints)
24. [🧪 Testing Endpoints](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%A7%AA-testing-endpoints)
25. [Pakai cURL (Command Line)](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-pakai-curl-command-line)
26. [Pakai GUI Tool](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-pakai-gui-tool)
27. [🏗️ Build Binary](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%8F%97%EF%B8%8F-build-binary)
28. [Build Standar](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-build-standar)
29. [Build Production (Smaller)](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-build-production-smaller)
30. [Cross-Compilation](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-cross-compilation)
31. [Jalankan Binary](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-jalankan-binary)
32. [☁️ Deploy ke Railway (Gratis!)](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%E2%98%81%EF%B8%8F-deploy-ke-railway-gratis)
33. [Persiapan](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-persiapan)
34. [Deploy di Railway](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-deploy-di-railway)
35. [🎓 Apa yang Sudah Kita Pelajari?](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%8E%93-apa-yang-sudah-kita-pelajari)
36. [💡 Tips & Best Practices](https://docs.kodingworks.io/s/01e57b74-74e6-44df-ac02-7e30a2478528#h-%F0%9F%92%A1-tips-and-best-practices)

[Outline](https://www.getoutline.com/?ref=sharelink)