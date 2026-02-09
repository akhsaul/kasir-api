package memory

import (
	"errors"
	"testing"

	model "kasir-api/models"
)

func TestNewCategoryRepository(t *testing.T) {
	repo := NewCategoryRepository()

	if repo == nil {
		t.Error("NewCategoryRepository should return a non-nil repository")
	}
	if repo.categories == nil {
		t.Error("NewCategoryRepository should initialize categories map")
	}
	if repo.nextCategoryID != 1 {
		t.Errorf("NewCategoryRepository should set nextCategoryID to 1, got: %d", repo.nextCategoryID)
	}
}

func TestCategoryRepository_GetAll_Empty(t *testing.T) {
	repo := NewCategoryRepository()

	categories, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("GetAll should return empty slice, got: %d", len(categories))
	}
}

func TestCategoryRepository_GetAll_WithData(t *testing.T) {
	repo := NewCategoryRepository()

	repo.Create(&model.Category{Name: "Electronics", Description: "Electronic items"})
	repo.Create(&model.Category{Name: "Food", Description: "Food items"})

	categories, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("GetAll should return 2 categories, got: %d", len(categories))
	}
}

func TestCategoryRepository_GetByID_Success(t *testing.T) {
	repo := NewCategoryRepository()

	created := &model.Category{Name: "Electronics", Description: "Electronic items"}
	repo.Create(created)

	category, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("GetByID should not return error, got: %v", err)
	}
	if category.ID != 1 {
		t.Errorf("GetByID should return category with ID 1, got: %d", category.ID)
	}
	if category.Name != "Electronics" {
		t.Errorf("GetByID should return category with name Electronics, got: %s", category.Name)
	}
}

func TestCategoryRepository_GetByID_NotFound(t *testing.T) {
	repo := NewCategoryRepository()

	_, err := repo.GetByID(999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("GetByID should return ErrNotFound, got: %v", err)
	}
}

func TestCategoryRepository_Create_Success(t *testing.T) {
	repo := NewCategoryRepository()

	category := &model.Category{Name: "Electronics", Description: "Electronic items"}
	err := repo.Create(category)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if category.ID != 1 {
		t.Errorf("Create should set ID to 1, got: %d", category.ID)
	}
	if repo.nextCategoryID != 2 {
		t.Errorf("Create should increment nextCategoryID to 2, got: %d", repo.nextCategoryID)
	}
}

func TestCategoryRepository_Create_Multiple(t *testing.T) {
	repo := NewCategoryRepository()

	cat1 := &model.Category{Name: "Electronics", Description: "Electronic items"}
	cat2 := &model.Category{Name: "Food", Description: "Food items"}

	repo.Create(cat1)
	repo.Create(cat2)

	if cat1.ID != 1 {
		t.Errorf("First category should have ID 1, got: %d", cat1.ID)
	}
	if cat2.ID != 2 {
		t.Errorf("Second category should have ID 2, got: %d", cat2.ID)
	}
	if repo.nextCategoryID != 3 {
		t.Errorf("nextCategoryID should be 3, got: %d", repo.nextCategoryID)
	}
}

func TestCategoryRepository_Update_Success(t *testing.T) {
	repo := NewCategoryRepository()

	category := &model.Category{Name: "Electronics", Description: "Electronic items"}
	repo.Create(category)

	category.Name = "Updated Electronics"
	category.Description = "Updated description"
	err := repo.Update(category)
	if err != nil {
		t.Errorf("Update should not return error, got: %v", err)
	}

	updated, _ := repo.GetByID(1)
	if updated.Name != "Updated Electronics" {
		t.Errorf("Update should change name, got: %s", updated.Name)
	}
	if updated.Description != "Updated description" {
		t.Errorf("Update should change description, got: %s", updated.Description)
	}
}

func TestCategoryRepository_Update_NotFound(t *testing.T) {
	repo := NewCategoryRepository()

	category := &model.Category{ID: 999, Name: "Test", Description: "Test"}
	err := repo.Update(category)

	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Update should return ErrNotFound, got: %v", err)
	}
}

func TestCategoryRepository_Delete_Success(t *testing.T) {
	repo := NewCategoryRepository()

	category := &model.Category{Name: "Electronics", Description: "Electronic items"}
	repo.Create(category)

	err := repo.Delete(1)
	if err != nil {
		t.Errorf("Delete should not return error, got: %v", err)
	}

	_, err = repo.GetByID(1)
	if !errors.Is(err, model.ErrNotFound) {
		t.Error("Delete should remove category from storage")
	}
}

func TestCategoryRepository_Delete_NotFound(t *testing.T) {
	repo := NewCategoryRepository()

	err := repo.Delete(999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Delete should return ErrNotFound, got: %v", err)
	}
}

func TestCategoryRepository_Concurrency(t *testing.T) {
	repo := NewCategoryRepository()

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			category := &model.Category{Name: "Category", Description: "Description"}
			repo.Create(category)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	categories, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(categories) != 10 {
		t.Errorf("Should have 10 categories, got: %d", len(categories))
	}
}

func TestCategoryRepository_GetAll_ReturnsCopy(t *testing.T) {
	repo := NewCategoryRepository()

	category := &model.Category{Name: "Electronics", Description: "Electronic items"}
	repo.Create(category)

	categories, _ := repo.GetAll()

	// Verify we get the same data
	if categories[0].Name != "Electronics" {
		t.Errorf("Category name should be Electronics, got: %s", categories[0].Name)
	}
}
