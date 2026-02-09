package memory

import (
	"errors"
	"testing"

	model "kasir-api/models"
)

// MockCategoryRepo is a simple mock for category repository
type MockCategoryRepo struct {
	categories map[int]*model.Category
}

func NewMockCategoryRepo() *MockCategoryRepo {
	return &MockCategoryRepo{
		categories: make(map[int]*model.Category),
	}
}

func (m *MockCategoryRepo) GetAll() ([]*model.Category, error) {
	categories := make([]*model.Category, 0, len(m.categories))
	for _, c := range m.categories {
		categories = append(categories, c)
	}
	return categories, nil
}

func (m *MockCategoryRepo) GetByID(id int) (*model.Category, error) {
	c, exists := m.categories[id]
	if !exists {
		return nil, model.ErrNotFound
	}
	return c, nil
}

func (m *MockCategoryRepo) Create(category *model.Category) error {
	category.ID = len(m.categories) + 1
	m.categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepo) Update(category *model.Category) error {
	if _, exists := m.categories[category.ID]; !exists {
		return model.ErrNotFound
	}
	m.categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepo) Delete(id int) error {
	if _, exists := m.categories[id]; !exists {
		return model.ErrNotFound
	}
	delete(m.categories, id)
	return nil
}

func TestNewProductRepository(t *testing.T) {
	repo := NewProductRepository(nil)

	if repo == nil {
		t.Error("NewProductRepository should return a non-nil repository")
	}
	if repo.products == nil {
		t.Error("NewProductRepository should initialize products map")
	}
	if repo.nextProductID != 1 {
		t.Errorf("NewProductRepository should set nextProductID to 1, got: %d", repo.nextProductID)
	}
}

func TestNewProductRepository_WithCategoryRepo(t *testing.T) {
	categoryRepo := NewMockCategoryRepo()
	repo := NewProductRepository(categoryRepo)

	if repo.categoryRepo != categoryRepo {
		t.Error("NewProductRepository should set categoryRepo")
	}
}

func TestProductRepository_GetAll_Empty(t *testing.T) {
	repo := NewProductRepository(nil)

	products, err := repo.GetAll("")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("GetAll should return empty slice, got: %d", len(products))
	}
}

func TestProductRepository_GetAll_WithData(t *testing.T) {
	repo := NewProductRepository(nil)

	repo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})
	repo.Create(&model.Product{Name: "Phone", Price: 500, Stock: 20})

	products, err := repo.GetAll("")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("GetAll should return 2 products, got: %d", len(products))
	}
}

func TestProductRepository_GetAll_WithNameFilter(t *testing.T) {
	repo := NewProductRepository(nil)

	repo.Create(&model.Product{Name: "Laptop Pro", Price: 1500, Stock: 10})
	repo.Create(&model.Product{Name: "Laptop Basic", Price: 1000, Stock: 15})
	repo.Create(&model.Product{Name: "Phone", Price: 500, Stock: 20})

	products, err := repo.GetAll("laptop")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("GetAll with filter 'laptop' should return 2 products, got: %d", len(products))
	}
}

func TestProductRepository_GetAll_CaseInsensitiveFilter(t *testing.T) {
	repo := NewProductRepository(nil)

	repo.Create(&model.Product{Name: "LAPTOP", Price: 1000, Stock: 10})

	products, err := repo.GetAll("laptop")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("GetAll filter should be case-insensitive, got: %d products", len(products))
	}
}

func TestProductRepository_GetAll_NoMatch(t *testing.T) {
	repo := NewProductRepository(nil)

	repo.Create(&model.Product{Name: "Laptop", Price: 1000, Stock: 10})

	products, err := repo.GetAll("phone")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("GetAll with no match should return 0 products, got: %d", len(products))
	}
}

func TestProductRepository_GetByID_Success(t *testing.T) {
	repo := NewProductRepository(nil)

	created := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	repo.Create(created)

	product, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("GetByID should not return error, got: %v", err)
	}
	if product.ID != 1 {
		t.Errorf("GetByID should return product with ID 1, got: %d", product.ID)
	}
	if product.Name != "Laptop" {
		t.Errorf("GetByID should return product with name Laptop, got: %s", product.Name)
	}
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	repo := NewProductRepository(nil)

	_, err := repo.GetByID(999)
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("GetByID should return ErrProductNotFound, got: %v", err)
	}
}

func TestProductRepository_GetByID_WithCategory(t *testing.T) {
	categoryRepo := NewMockCategoryRepo()
	categoryRepo.categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	repo := NewProductRepository(categoryRepo)

	categoryID := 1
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	repo.Create(product)

	retrieved, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("GetByID should not return error, got: %v", err)
	}
	if retrieved.Category == nil {
		t.Error("GetByID should enrich product with category")
	}
	if retrieved.Category != nil && retrieved.Category.Name != "Electronics" {
		t.Errorf("Category name should be Electronics, got: %s", retrieved.Category.Name)
	}
}

func TestProductRepository_Create_Success(t *testing.T) {
	repo := NewProductRepository(nil)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	err := repo.Create(product)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if product.ID != 1 {
		t.Errorf("Create should set ID to 1, got: %d", product.ID)
	}
	if repo.nextProductID != 2 {
		t.Errorf("Create should increment nextProductID to 2, got: %d", repo.nextProductID)
	}
}

func TestProductRepository_Create_Multiple(t *testing.T) {
	repo := NewProductRepository(nil)

	prod1 := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	prod2 := &model.Product{Name: "Phone", Price: 500, Stock: 20}

	repo.Create(prod1)
	repo.Create(prod2)

	if prod1.ID != 1 {
		t.Errorf("First product should have ID 1, got: %d", prod1.ID)
	}
	if prod2.ID != 2 {
		t.Errorf("Second product should have ID 2, got: %d", prod2.ID)
	}
	if repo.nextProductID != 3 {
		t.Errorf("nextProductID should be 3, got: %d", repo.nextProductID)
	}
}

func TestProductRepository_Create_WithCategory(t *testing.T) {
	categoryRepo := NewMockCategoryRepo()
	categoryRepo.categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	repo := NewProductRepository(categoryRepo)

	categoryID := 1
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	err := repo.Create(product)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if product.Category == nil {
		t.Error("Create should enrich product with category")
	}
}

func TestProductRepository_Update_Success(t *testing.T) {
	repo := NewProductRepository(nil)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	repo.Create(product)

	product.Name = "Updated Laptop"
	product.Price = 1500
	err := repo.Update(product)
	if err != nil {
		t.Errorf("Update should not return error, got: %v", err)
	}

	updated, _ := repo.GetByID(1)
	if updated.Name != "Updated Laptop" {
		t.Errorf("Update should change name, got: %s", updated.Name)
	}
	if updated.Price != 1500 {
		t.Errorf("Update should change price, got: %d", updated.Price)
	}
}

func TestProductRepository_Update_NotFound(t *testing.T) {
	repo := NewProductRepository(nil)

	product := &model.Product{ID: 999, Name: "Test", Price: 100, Stock: 10}
	err := repo.Update(product)

	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("Update should return ErrProductNotFound, got: %v", err)
	}
}

func TestProductRepository_Update_WithCategory(t *testing.T) {
	categoryRepo := NewMockCategoryRepo()
	categoryRepo.categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	repo := NewProductRepository(categoryRepo)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	repo.Create(product)

	categoryID := 1
	product.CategoryID = &categoryID
	product.Name = "Updated Laptop"
	err := repo.Update(product)
	if err != nil {
		t.Errorf("Update should not return error, got: %v", err)
	}
	if product.Category == nil {
		t.Error("Update should enrich product with category")
	}
}

func TestProductRepository_Delete_Success(t *testing.T) {
	repo := NewProductRepository(nil)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	repo.Create(product)

	err := repo.Delete(1)
	if err != nil {
		t.Errorf("Delete should not return error, got: %v", err)
	}

	_, err = repo.GetByID(1)
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Error("Delete should remove product from storage")
	}
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	repo := NewProductRepository(nil)

	err := repo.Delete(999)
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("Delete should return ErrProductNotFound, got: %v", err)
	}
}

func TestProductRepository_Concurrency(t *testing.T) {
	repo := NewProductRepository(nil)

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			product := &model.Product{Name: "Product", Price: 100, Stock: 10}
			repo.Create(product)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	products, err := repo.GetAll("")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 10 {
		t.Errorf("Should have 10 products, got: %d", len(products))
	}
}

func TestProductRepository_GetAll_ReturnsCopy(t *testing.T) {
	repo := NewProductRepository(nil)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	repo.Create(product)

	products, _ := repo.GetAll("")

	// Modify the returned product
	products[0].Name = "Modified"

	// Original should not be affected
	original, _ := repo.GetByID(1)
	if original.Name == "Modified" {
		t.Error("GetAll should return copies, not references")
	}
}

func TestProductRepository_GetByID_ReturnsCopy(t *testing.T) {
	repo := NewProductRepository(nil)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	repo.Create(product)

	retrieved, _ := repo.GetByID(1)
	retrieved.Name = "Modified"

	// Original should not be affected
	original, _ := repo.GetByID(1)
	if original.Name == "Modified" {
		t.Error("GetByID should return a copy, not reference")
	}
}

func TestProductRepository_EnrichWithCategory_NilCategoryRepo(t *testing.T) {
	repo := NewProductRepository(nil)

	categoryID := 1
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	repo.Create(product)

	retrieved, _ := repo.GetByID(1)
	if retrieved.Category != nil {
		t.Error("Product category should be nil when categoryRepo is nil")
	}
}

func TestProductRepository_EnrichWithCategory_NilCategoryID(t *testing.T) {
	categoryRepo := NewMockCategoryRepo()
	repo := NewProductRepository(categoryRepo)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: nil}
	repo.Create(product)

	retrieved, _ := repo.GetByID(1)
	if retrieved.Category != nil {
		t.Error("Product category should be nil when CategoryID is nil")
	}
}

func TestProductRepository_EnrichWithCategory_CategoryNotFound(t *testing.T) {
	categoryRepo := NewMockCategoryRepo()
	repo := NewProductRepository(categoryRepo)

	categoryID := 999
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	repo.Create(product)

	retrieved, _ := repo.GetByID(1)
	if retrieved.Category != nil {
		t.Error("Product category should be nil when category not found")
	}
}
