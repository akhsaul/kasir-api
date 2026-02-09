package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"kasir-api/config"
	"kasir-api/helpers/logger"
	"kasir-api/repositories/postgres"
)

// Seed data
var categories = []struct {
	Name        string
	Description string
}{
	{"Makanan", "Makanan instan, snack, dan cemilan"},
	{"Minuman", "Minuman kemasan dan botol"},
	{"Rokok", "Rokok dan produk tembakau"},
	{"Kebutuhan Rumah", "Sabun, deterjen, dan kebutuhan rumah tangga"},
	{"Alat Tulis", "Pulpen, buku, dan alat tulis kantor"},
}

var products = []struct {
	Name       string
	Price      int
	Stock      int
	CategoryID int
}{
	// Makanan (category_id: 1)
	{"Indomie Goreng", 3500, 100, 1},
	{"Indomie Rebus", 3000, 80, 1},
	{"Mie Sedaap Goreng", 3500, 75, 1},
	{"Pop Mie Ayam", 6000, 50, 1},
	{"Chitato Original", 12000, 40, 1},
	{"Lays Classic", 13000, 35, 1},
	{"Oreo Original", 8500, 45, 1},
	{"Biskuit Roma Kelapa", 5000, 60, 1},
	{"Roti Tawar Sari Roti", 15000, 25, 1},
	{"Energen Coklat", 2500, 90, 1},

	// Minuman (category_id: 2)
	{"Aqua 600ml", 4000, 120, 2},
	{"Aqua 1500ml", 7500, 60, 2},
	{"Teh Botol Sosro", 5000, 80, 2},
	{"Teh Pucuk Harum", 4500, 85, 2},
	{"Coca Cola 390ml", 7000, 50, 2},
	{"Fanta Strawberry 390ml", 7000, 45, 2},
	{"Sprite 390ml", 7000, 40, 2},
	{"Pocari Sweat 500ml", 8500, 55, 2},
	{"Kopi Good Day Cappucino", 4000, 70, 2},
	{"Ultra Milk Coklat 250ml", 6000, 65, 2},

	// Rokok (category_id: 3)
	{"Gudang Garam Surya 16", 28000, 30, 3},
	{"Djarum Super 16", 26000, 35, 3},
	{"Sampoerna Mild 16", 32000, 40, 3},
	{"Marlboro Red 20", 38000, 25, 3},
	{"LA Lights 16", 24000, 45, 3},

	// Kebutuhan Rumah (category_id: 4)
	{"Rinso Cair 800ml", 18000, 30, 4},
	{"Sunlight Jeruk 755ml", 15000, 35, 4},
	{"Molto Pewangi 800ml", 16000, 28, 4},
	{"Baygon Aerosol 600ml", 45000, 20, 4},
	{"Tissue Paseo 250 sheets", 18000, 40, 4},
	{"Sabun Lifebuoy 85g", 4000, 60, 4},
	{"Shampo Pantene 170ml", 28000, 25, 4},
	{"Pasta Gigi Pepsodent 190g", 15000, 35, 4},

	// Alat Tulis (category_id: 5)
	{"Pulpen Standard AE7", 3000, 100, 5},
	{"Pensil 2B Faber Castell", 2500, 80, 5},
	{"Buku Tulis Sidu 58 lembar", 5000, 50, 5},
	{"Penghapus Steadler", 3500, 60, 5},
	{"Penggaris 30cm", 4000, 40, 5},
	{"Tip-X Kenko", 8000, 35, 5},
}

func main() {
	logger.InitLogger()

	// Parse flags
	seedTransactions := flag.Bool("transactions", false, "Also seed sample transactions")
	clearData := flag.Bool("clear", false, "Clear existing data before seeding")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	if !cfg.DB.Enabled {
		logger.Fatal("Database is not enabled. Please set DB_ENABLED=true in .env")
	}

	db, err := postgres.NewDB(&cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("Starting database seeder...")

	if *clearData {
		if err := clearAllData(db); err != nil {
			logger.Fatal(err)
		}
		logger.Info("Cleared existing data")
	}

	// Seed categories
	categoryIDs, err := seedCategories(db)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Seeded %d categories", len(categoryIDs))

	// Seed products
	productIDs, err := seedProducts(db)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Seeded %d products", len(productIDs))

	// Seed transactions if flag is set
	if *seedTransactions {
		count, err := seedSampleTransactions(db, productIDs)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Seeded %d transactions", count)
	}

	logger.Info("Seeding completed successfully!")
}

func clearAllData(db *postgres.DB) error {
	queries := []string{
		"DELETE FROM transaction_details",
		"DELETE FROM transactions",
		"DELETE FROM products",
		"DELETE FROM categories",
		"ALTER SEQUENCE IF EXISTS categories_id_seq RESTART WITH 1",
		"ALTER SEQUENCE IF EXISTS products_id_seq RESTART WITH 1",
		"ALTER SEQUENCE IF EXISTS transactions_id_seq RESTART WITH 1",
		"ALTER SEQUENCE IF EXISTS transaction_details_id_seq RESTART WITH 1",
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("clear data: %w", err)
		}
	}
	return nil
}

func seedCategories(db *postgres.DB) ([]int, error) {
	var ids []int
	for _, c := range categories {
		var id int
		err := db.QueryRow(`
			INSERT INTO categories (name, description) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING
			RETURNING id
		`, c.Name, c.Description).Scan(&id)
		if err != nil {
			// Try to get existing ID
			err = db.QueryRow("SELECT id FROM categories WHERE name = $1", c.Name).Scan(&id)
			if err != nil {
				return nil, fmt.Errorf("seed category %s: %w", c.Name, err)
			}
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedProducts(db *postgres.DB) ([]int, error) {
	var ids []int
	for _, p := range products {
		var id int
		err := db.QueryRow(`
			INSERT INTO products (name, price, stock, category_id) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT DO NOTHING
			RETURNING id
		`, p.Name, p.Price, p.Stock, p.CategoryID).Scan(&id)
		if err != nil {
			// Try to get existing ID
			err = db.QueryRow("SELECT id FROM products WHERE name = $1", p.Name).Scan(&id)
			if err != nil {
				return nil, fmt.Errorf("seed product %s: %w", p.Name, err)
			}
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedSampleTransactions(db *postgres.DB, productIDs []int) (int, error) {
	if len(productIDs) == 0 {
		return 0, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	transactionCount := 10

	for i := 0; i < transactionCount; i++ {
		// Random number of items per transaction (1-5)
		itemCount := r.Intn(5) + 1

		// Create transaction
		var transactionID int
		createdAt := time.Now().Add(-time.Duration(r.Intn(7*24)) * time.Hour) // Random within last 7 days

		err := db.QueryRow(`
			INSERT INTO transactions (total_amount, created_at) 
			VALUES (0, $1) 
			RETURNING id
		`, createdAt).Scan(&transactionID)
		if err != nil {
			return 0, fmt.Errorf("create transaction: %w", err)
		}

		totalAmount := 0
		usedProducts := make(map[int]bool)

		for j := 0; j < itemCount; j++ {
			// Pick random product (avoid duplicates in same transaction)
			var productID int
			for {
				productID = productIDs[r.Intn(len(productIDs))]
				if !usedProducts[productID] {
					usedProducts[productID] = true
					break
				}
				if len(usedProducts) >= len(productIDs) {
					break
				}
			}

			// Get product info
			var productName string
			var price int
			err := db.QueryRow("SELECT name, price FROM products WHERE id = $1", productID).Scan(&productName, &price)
			if err != nil {
				continue
			}

			quantity := r.Intn(5) + 1
			subtotal := price * quantity
			totalAmount += subtotal

			// Insert transaction detail
			_, err = db.Exec(`
				INSERT INTO transaction_details (transaction_id, product_id, product_name, quantity, price, subtotal)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, transactionID, productID, productName, quantity, price, subtotal)
			if err != nil {
				return 0, fmt.Errorf("create transaction detail: %w", err)
			}
		}

		// Update total amount
		_, err = db.Exec("UPDATE transactions SET total_amount = $1 WHERE id = $2", totalAmount, transactionID)
		if err != nil {
			return 0, fmt.Errorf("update transaction total: %w", err)
		}
	}

	return transactionCount, nil
}
