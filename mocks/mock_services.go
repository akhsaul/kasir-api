package mocks

import (
	"time"

	model "kasir-api/models"
)

// MockCategoryService is a mock implementation for testing handlers.
type MockCategoryService struct {
	Categories  map[int]*model.Category
	GetAllFunc  func() ([]*model.Category, error)
	GetByIDFunc func(id int) (*model.Category, error)
	CreateFunc  func(category *model.Category) (*model.Category, error)
	UpdateFunc  func(id int, category *model.Category) (*model.Category, error)
	DeleteFunc  func(id int) error
}

func NewMockCategoryService() *MockCategoryService {
	return &MockCategoryService{
		Categories: make(map[int]*model.Category),
	}
}

func (m *MockCategoryService) GetAll() ([]*model.Category, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	categories := make([]*model.Category, 0, len(m.Categories))
	for _, c := range m.Categories {
		categories = append(categories, c)
	}
	return categories, nil
}

func (m *MockCategoryService) GetByID(id int) (*model.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	c, exists := m.Categories[id]
	if !exists {
		return nil, model.ErrCategoryNotFound
	}
	return c, nil
}

func (m *MockCategoryService) Create(category *model.Category) (*model.Category, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(category)
	}
	category.ID = len(m.Categories) + 1
	m.Categories[category.ID] = category
	return category, nil
}

func (m *MockCategoryService) Update(id int, category *model.Category) (*model.Category, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, category)
	}
	if _, exists := m.Categories[id]; !exists {
		return nil, model.ErrCategoryNotFound
	}
	category.ID = id
	m.Categories[id] = category
	return category, nil
}

func (m *MockCategoryService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if _, exists := m.Categories[id]; !exists {
		return model.ErrCategoryNotFound
	}
	delete(m.Categories, id)
	return nil
}

// MockProductService is a mock implementation for testing handlers.
type MockProductService struct {
	Products    map[int]*model.Product
	GetAllFunc  func(name string) ([]*model.Product, error)
	GetByIDFunc func(id int) (*model.Product, error)
	CreateFunc  func(product *model.Product) (*model.Product, error)
	UpdateFunc  func(id int, product *model.Product) (*model.Product, error)
	DeleteFunc  func(id int) error
}

func NewMockProductService() *MockProductService {
	return &MockProductService{
		Products: make(map[int]*model.Product),
	}
}

func (m *MockProductService) GetAll(name string) ([]*model.Product, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(name)
	}
	products := make([]*model.Product, 0, len(m.Products))
	for _, p := range m.Products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MockProductService) GetByID(id int) (*model.Product, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	p, exists := m.Products[id]
	if !exists {
		return nil, model.ErrProductNotFound
	}
	return p, nil
}

func (m *MockProductService) Create(product *model.Product) (*model.Product, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(product)
	}
	product.ID = len(m.Products) + 1
	m.Products[product.ID] = product
	return product, nil
}

func (m *MockProductService) Update(id int, product *model.Product) (*model.Product, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, product)
	}
	if _, exists := m.Products[id]; !exists {
		return nil, model.ErrProductNotFound
	}
	product.ID = id
	m.Products[id] = product
	return product, nil
}

func (m *MockProductService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if _, exists := m.Products[id]; !exists {
		return model.ErrProductNotFound
	}
	delete(m.Products, id)
	return nil
}

// MockTransactionService is a mock implementation for testing handlers.
type MockTransactionService struct {
	Transactions             map[int]*model.Transaction
	CheckoutFunc             func(request *model.CheckoutRequest) (*model.Transaction, error)
	GetByIDFunc              func(id int) (*model.Transaction, error)
	GetTodayReportFunc       func() (*model.ReportResponse, error)
	GetReportByDateRangeFunc func(startDate, endDate time.Time) (*model.ReportResponse, error)
}

func NewMockTransactionService() *MockTransactionService {
	return &MockTransactionService{
		Transactions: make(map[int]*model.Transaction),
	}
}

func (m *MockTransactionService) Checkout(request *model.CheckoutRequest) (*model.Transaction, error) {
	if m.CheckoutFunc != nil {
		return m.CheckoutFunc(request)
	}
	transaction := &model.Transaction{
		ID:          len(m.Transactions) + 1,
		TotalAmount: 1000,
		CreatedAt:   time.Now(),
	}
	m.Transactions[transaction.ID] = transaction
	return transaction, nil
}

func (m *MockTransactionService) GetByID(id int) (*model.Transaction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	t, exists := m.Transactions[id]
	if !exists {
		return nil, model.ErrTransactionNotFound
	}
	return t, nil
}

func (m *MockTransactionService) GetTodayReport() (*model.ReportResponse, error) {
	if m.GetTodayReportFunc != nil {
		return m.GetTodayReportFunc()
	}
	return &model.ReportResponse{}, nil
}

func (m *MockTransactionService) GetReportByDateRange(startDate, endDate time.Time) (*model.ReportResponse, error) {
	if m.GetReportByDateRangeFunc != nil {
		return m.GetReportByDateRangeFunc(startDate, endDate)
	}
	return &model.ReportResponse{}, nil
}
