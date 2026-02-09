package service

import (
	"errors"
	"testing"

	"kasir-api/mocks"
	model "kasir-api/models"
)

func TestNewCategoryService(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	if service == nil {
		t.Error("NewCategoryService should return a non-nil service")
	}
	if service.repo != repo {
		t.Error("NewCategoryService should set the repo")
	}
}

func TestCategoryService_GetAll(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	// Add test data
	repo.Categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}
	repo.Categories[2] = &model.Category{ID: 2, Name: "Food", Description: "Food items"}

	categories, err := service.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error, got: %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("GetAll should return 2 categories, got: %d", len(categories))
	}
}

func TestCategoryService_GetAll_Error(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	expectedErr := errors.New("database error")
	repo.GetAllFunc = func() ([]*model.Category, error) {
		return nil, expectedErr
	}

	_, err := service.GetAll()
	if err != expectedErr {
		t.Errorf("GetAll should return the error from repo, got: %v", err)
	}
}

func TestCategoryService_GetByID_Success(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	repo.Categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	category, err := service.GetByID(1)
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

func TestCategoryService_GetByID_InvalidID(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

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
			if !errors.Is(err, model.ErrIDRequired) {
				t.Errorf("GetByID with %s id should return ErrIDRequired, got: %v", tc.name, err)
			}
		})
	}
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	_, err := service.GetByID(999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("GetByID should return ErrNotFound, got: %v", err)
	}
}

func TestCategoryService_Create_Success(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	category := &model.Category{Name: "Electronics", Description: "Electronic items"}
	created, err := service.Create(category)
	if err != nil {
		t.Errorf("Create should not return error, got: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("Create should return category with ID 1, got: %d", created.ID)
	}
	if created.Name != "Electronics" {
		t.Errorf("Create should return category with name Electronics, got: %s", created.Name)
	}
}

func TestCategoryService_Create_EmptyName(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	testCases := []struct {
		name     string
		category *model.Category
	}{
		{"empty name", &model.Category{Name: "", Description: "desc"}},
		{"whitespace name", &model.Category{Name: "   ", Description: "desc"}},
		{"tab only", &model.Category{Name: "\t", Description: "desc"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Create(tc.category)
			if !errors.Is(err, model.ErrNameRequired) {
				t.Errorf("Create with %s should return ErrNameRequired, got: %v", tc.name, err)
			}
		})
	}
}

func TestCategoryService_Create_RepoError(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	expectedErr := errors.New("database error")
	repo.CreateFunc = func(category *model.Category) error {
		return expectedErr
	}

	category := &model.Category{Name: "Electronics", Description: "Electronic items"}
	_, err := service.Create(category)
	if err != expectedErr {
		t.Errorf("Create should return the error from repo, got: %v", err)
	}
}

func TestCategoryService_Update_Success(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	repo.Categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	updated, err := service.Update(1, &model.Category{Name: "Updated Electronics", Description: "Updated desc"})
	if err != nil {
		t.Errorf("Update should not return error, got: %v", err)
	}
	if updated.ID != 1 {
		t.Errorf("Update should return category with ID 1, got: %d", updated.ID)
	}
	if updated.Name != "Updated Electronics" {
		t.Errorf("Update should return category with updated name, got: %s", updated.Name)
	}
}

func TestCategoryService_Update_InvalidID(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	testCases := []struct {
		name string
		id   int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Update(tc.id, &model.Category{Name: "Test"})
			if !errors.Is(err, model.ErrIDRequired) {
				t.Errorf("Update with %s id should return ErrIDRequired, got: %v", tc.name, err)
			}
		})
	}
}

func TestCategoryService_Update_EmptyName(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	repo.Categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	_, err := service.Update(1, &model.Category{Name: "", Description: "desc"})
	if !errors.Is(err, model.ErrNameRequired) {
		t.Errorf("Update with empty name should return ErrNameRequired, got: %v", err)
	}
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	_, err := service.Update(999, &model.Category{Name: "Test"})
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Update should return ErrNotFound, got: %v", err)
	}
}

func TestCategoryService_Update_RepoError(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	expectedErr := errors.New("database error")
	repo.UpdateFunc = func(category *model.Category) error {
		return expectedErr
	}

	_, err := service.Update(1, &model.Category{Name: "Test"})
	if err != expectedErr {
		t.Errorf("Update should return the error from repo, got: %v", err)
	}
}

func TestCategoryService_Delete_Success(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	repo.Categories[1] = &model.Category{ID: 1, Name: "Electronics", Description: "Electronic items"}

	err := service.Delete(1)
	if err != nil {
		t.Errorf("Delete should not return error, got: %v", err)
	}
	if _, exists := repo.Categories[1]; exists {
		t.Error("Delete should remove the category from repo")
	}
}

func TestCategoryService_Delete_InvalidID(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

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
			if !errors.Is(err, model.ErrIDRequired) {
				t.Errorf("Delete with %s id should return ErrIDRequired, got: %v", tc.name, err)
			}
		})
	}
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	repo := mocks.NewMockCategoryRepository()
	service := NewCategoryService(repo)

	err := service.Delete(999)
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("Delete should return ErrNotFound, got: %v", err)
	}
}
