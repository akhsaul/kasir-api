package mocks

import (
	"time"

	model "kasir-api/models"
)

// MockCategoryRepository is a mock implementation of repository.CategoryRepository.
type MockCategoryRepository struct {
	Categories  map[int]*model.Category
	NextID      int
	GetAllFunc  func() ([]*model.Category, error)
	GetByIDFunc func(id int) (*model.Category, error)
	CreateFunc  func(category *model.Category) error
	UpdateFunc  func(category *model.Category) error
	DeleteFunc  func(id int) error
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		Categories: make(map[int]*model.Category),
		NextID:     1,
	}
}

func (m *MockCategoryRepository) GetAll() ([]*model.Category, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	categories := make([]*model.Category, 0, len(m.Categories))
	for _, c := range m.Categories {
		categories = append(categories, c)
	}
	return categories, nil
}

func (m *MockCategoryRepository) GetByID(id int) (*model.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	c, exists := m.Categories[id]
	if !exists {
		return nil, model.ErrCategoryNotFound
	}
	return c, nil
}

func (m *MockCategoryRepository) Create(category *model.Category) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(category)
	}
	category.ID = m.NextID
	m.Categories[category.ID] = category
	m.NextID++
	return nil
}

func (m *MockCategoryRepository) Update(category *model.Category) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(category)
	}
	if _, exists := m.Categories[category.ID]; !exists {
		return model.ErrCategoryNotFound
	}
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if _, exists := m.Categories[id]; !exists {
		return model.ErrCategoryNotFound
	}
	delete(m.Categories, id)
	return nil
}

// MockProductRepository is a mock implementation of repository.ProductRepository.
type MockProductRepository struct {
	Products    map[int]*model.Product
	NextID      int
	GetAllFunc  func(name string) ([]*model.Product, error)
	GetByIDFunc func(id int) (*model.Product, error)
	CreateFunc  func(product *model.Product) error
	UpdateFunc  func(product *model.Product) error
	DeleteFunc  func(id int) error
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products: make(map[int]*model.Product),
		NextID:   1,
	}
}

func (m *MockProductRepository) GetAll(name string) ([]*model.Product, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(name)
	}
	products := make([]*model.Product, 0, len(m.Products))
	for _, p := range m.Products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MockProductRepository) GetByID(id int) (*model.Product, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	p, exists := m.Products[id]
	if !exists {
		return nil, model.ErrProductNotFound
	}
	return p, nil
}

func (m *MockProductRepository) Create(product *model.Product) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(product)
	}
	product.ID = m.NextID
	m.Products[product.ID] = product
	m.NextID++
	return nil
}

func (m *MockProductRepository) Update(product *model.Product) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(product)
	}
	if _, exists := m.Products[product.ID]; !exists {
		return model.ErrProductNotFound
	}
	m.Products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if _, exists := m.Products[id]; !exists {
		return model.ErrProductNotFound
	}
	delete(m.Products, id)
	return nil
}

// MockTransactionRepository is a mock implementation of repository.TransactionRepository.
type MockTransactionRepository struct {
	Transactions             map[int]*model.Transaction
	NextID                   int
	CreateFunc               func(transaction *model.Transaction) error
	GetByIDFunc              func(id int) (*model.Transaction, error)
	GetReportByDateRangeFunc func(startDate, endDate time.Time) (*model.ReportResponse, error)
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		Transactions: make(map[int]*model.Transaction),
		NextID:       1,
	}
}

func (m *MockTransactionRepository) Create(transaction *model.Transaction) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(transaction)
	}
	transaction.ID = m.NextID
	m.Transactions[transaction.ID] = transaction
	m.NextID++
	return nil
}

func (m *MockTransactionRepository) GetByID(id int) (*model.Transaction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	t, exists := m.Transactions[id]
	if !exists {
		return nil, model.ErrTransactionNotFound
	}
	return t, nil
}

func (m *MockTransactionRepository) GetReportByDateRange(startDate, endDate time.Time) (*model.ReportResponse, error) {
	if m.GetReportByDateRangeFunc != nil {
		return m.GetReportByDateRangeFunc(startDate, endDate)
	}
	return &model.ReportResponse{}, nil
}
