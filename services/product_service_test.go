package service

import (
	"errors"
	"testing"

	"kasir-api/mocks"
	model "kasir-api/models"
)

func TestNewProductService(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	if service == nil {
		t.Error("NewProductService should return a non-nil service")
	}
	if service.repo != productRepo {
		t.Error("NewProductService should set the product repo")
	}
	if service.categoryRepo != categoryRepo {
		t.Error("NewProductService should set the category repo")
	}
}

func TestProductService_GetAll_Success(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}
	productRepo.Products[2] = &model.Product{ID: 2, Name: "Phone", Price: 500, Stock: 20}

	products, err := service.GetAll("")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("GetAll should return 2 products, got: %d", len(products))
	}
}

func TestProductService_GetAll_WithNameFilter(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.GetAllFunc = func(name string) ([]*model.Product, error) {
		if name == "laptop" {
			return []*model.Product{{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}}, nil
		}
		return []*model.Product{}, nil
	}

	products, err := service.GetAll("laptop")
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("GetAll should return 1 product, got: %d", len(products))
	}
}

func TestProductService_GetAll_Error(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	expectedErr := errors.New("database error")
	productRepo.GetAllFunc = func(name string) ([]*model.Product, error) {
		return nil, expectedErr
	}

	_, err := service.GetAll("")
	if err != expectedErr {
		t.Errorf("GetAll should return the error from repo, got: %v", err)
	}
}

func TestProductService_GetByID_Success(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	product, err := service.GetByID(1)
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

func TestProductService_GetByID_InvalidID(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	testCases := []struct {
		name string
		id   int
	}{
		{"zero", 0},
		{"negative", -1},
		{"negative large", -100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.GetByID(tc.id)
			if !errors.Is(err, model.ErrProductNotFound) {
				t.Errorf("GetByID with %s id should return ErrProductNotFound, got: %v", tc.name, err)
			}
		})
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	_, err := service.GetByID(999)
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("GetByID should return ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_Create_Success(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	created, err := service.Create(product)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("Create should return product with ID 1, got: %d", created.ID)
	}
	if created.Name != "Laptop" {
		t.Errorf("Create should return product with name Laptop, got: %s", created.Name)
	}
}

func TestProductService_Create_WithValidCategory(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	categoryRepo.Categories[1] = &model.Category{ID: 1, Name: "Electronics"}

	categoryID := 1
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	created, err := service.Create(product)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if created.CategoryID == nil || *created.CategoryID != 1 {
		t.Error("Create should return product with CategoryID 1")
	}
}

func TestProductService_Create_WithInvalidCategory(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	categoryID := 999
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	_, err := service.Create(product)

	if !errors.Is(err, model.ErrCategoryNotFound) {
		t.Errorf("Create with invalid category should return ErrCategoryNotFound, got: %v", err)
	}
}

func TestProductService_Create_CategoryRepoError(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	expectedErr := errors.New("database error")
	categoryRepo.GetByIDFunc = func(id int) (*model.Category, error) {
		return nil, expectedErr
	}

	categoryID := 1
	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10, CategoryID: &categoryID}
	_, err := service.Create(product)

	if err != expectedErr {
		t.Errorf("Create should return the error from category repo, got: %v", err)
	}
}

func TestProductService_Create_EmptyName(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	testCases := []struct {
		name    string
		product *model.Product
	}{
		{"empty name", &model.Product{Name: "", Price: 1000, Stock: 10}},
		{"whitespace name", &model.Product{Name: "   ", Price: 1000, Stock: 10}},
		{"tab only", &model.Product{Name: "\t", Price: 1000, Stock: 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Create(tc.product)
			if !errors.Is(err, model.ErrNameRequired) {
				t.Errorf("Create with %s should return ErrNameRequired, got: %v", tc.name, err)
			}
		})
	}
}

func TestProductService_Create_InvalidPrice(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	testCases := []struct {
		name    string
		product *model.Product
	}{
		{"zero price", &model.Product{Name: "Laptop", Price: 0, Stock: 10}},
		{"negative price", &model.Product{Name: "Laptop", Price: -100, Stock: 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Create(tc.product)
			if !errors.Is(err, model.ErrPriceInvalid) {
				t.Errorf("Create with %s should return ErrPriceInvalid, got: %v", tc.name, err)
			}
		})
	}
}

func TestProductService_Create_InvalidStock(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: -1}
	_, err := service.Create(product)

	if !errors.Is(err, model.ErrStockInvalid) {
		t.Errorf("Create with negative stock should return ErrStockInvalid, got: %v", err)
	}
}

func TestProductService_Create_RepoError(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	expectedErr := errors.New("database error")
	productRepo.CreateFunc = func(product *model.Product) error {
		return expectedErr
	}

	product := &model.Product{Name: "Laptop", Price: 1000, Stock: 10}
	_, err := service.Create(product)

	if err != expectedErr {
		t.Errorf("Create should return the error from repo, got: %v", err)
	}
}

func TestProductService_Update_Success(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	updated, err := service.Update(1, &model.Product{Name: "Updated Laptop", Price: 1500, Stock: 5})
	if err != nil {
		t.Errorf("Update should not return error, got: %v", err)
	}
	if updated.ID != 1 {
		t.Errorf("Update should return product with ID 1, got: %d", updated.ID)
	}
	if updated.Name != "Updated Laptop" {
		t.Errorf("Update should return product with updated name, got: %s", updated.Name)
	}
}

func TestProductService_Update_InvalidID(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	testCases := []struct {
		name string
		id   int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Update(tc.id, &model.Product{Name: "Test", Price: 100, Stock: 10})
			if !errors.Is(err, model.ErrProductNotFound) {
				t.Errorf("Update with %s id should return ErrProductNotFound, got: %v", tc.name, err)
			}
		})
	}
}

func TestProductService_Update_EmptyName(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	_, err := service.Update(1, &model.Product{Name: "", Price: 1000, Stock: 10})
	if !errors.Is(err, model.ErrNameRequired) {
		t.Errorf("Update with empty name should return ErrNameRequired, got: %v", err)
	}
}

func TestProductService_Update_InvalidPrice(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	_, err := service.Update(1, &model.Product{Name: "Laptop", Price: 0, Stock: 10})
	if !errors.Is(err, model.ErrPriceInvalid) {
		t.Errorf("Update with zero price should return ErrPriceInvalid, got: %v", err)
	}
}

func TestProductService_Update_InvalidStock(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	_, err := service.Update(1, &model.Product{Name: "Laptop", Price: 1000, Stock: -1})
	if !errors.Is(err, model.ErrStockInvalid) {
		t.Errorf("Update with negative stock should return ErrStockInvalid, got: %v", err)
	}
}

func TestProductService_Update_NotFound(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	_, err := service.Update(999, &model.Product{Name: "Test", Price: 100, Stock: 10})
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("Update should return ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_Update_RepoError(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	expectedErr := errors.New("database error")
	productRepo.UpdateFunc = func(product *model.Product) error {
		return expectedErr
	}

	_, err := service.Update(1, &model.Product{Name: "Test", Price: 100, Stock: 10})
	if err != expectedErr {
		t.Errorf("Update should return the error from repo, got: %v", err)
	}
}

func TestProductService_Delete_Success(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	productRepo.Products[1] = &model.Product{ID: 1, Name: "Laptop", Price: 1000, Stock: 10}

	err := service.Delete(1)
	if err != nil {
		t.Errorf("Delete should not return error, got: %v", err)
	}
	if _, exists := productRepo.Products[1]; exists {
		t.Error("Delete should remove the product from repo")
	}
}

func TestProductService_Delete_InvalidID(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	testCases := []struct {
		name string
		id   int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.Delete(tc.id)
			if !errors.Is(err, model.ErrProductNotFound) {
				t.Errorf("Delete with %s id should return ErrProductNotFound, got: %v", tc.name, err)
			}
		})
	}
}

func TestProductService_Delete_NotFound(t *testing.T) {
	productRepo := mocks.NewMockProductRepository()
	categoryRepo := mocks.NewMockCategoryRepository()
	service := NewProductService(productRepo, categoryRepo)

	err := service.Delete(999)
	if !errors.Is(err, model.ErrProductNotFound) {
		t.Errorf("Delete should return ErrProductNotFound, got: %v", err)
	}
}
