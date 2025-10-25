// Package repository demonstrates the Repository pattern.
// It mediates between domain and data layers, providing a collection-like
// interface for accessing domain objects while hiding data access details.
package repository

import (
	"errors"
	"fmt"
	"sync"
)

// User domain model
type User struct {
	ID    string
	Name  string
	Email string
	Age   int
}

// UserRepository interface defines data access operations
type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	Save(user *User) error
	Update(user *User) error
	Delete(id string) error
}

// InMemoryUserRepository implements UserRepository with in-memory storage
type InMemoryUserRepository struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewInMemoryUserRepository creates an in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) FindByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) FindAll() ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *InMemoryUserRepository) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Update(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

// Product domain model
type Product struct {
	ID       string
	Name     string
	Price    float64
	Category string
	InStock  bool
}

// ProductRepository interface
type ProductRepository interface {
	FindByID(id string) (*Product, error)
	FindByCategory(category string) ([]*Product, error)
	FindInStock() ([]*Product, error)
	Save(product *Product) error
	Update(product *Product) error
	Delete(id string) error
}

// InMemoryProductRepository implements ProductRepository
type InMemoryProductRepository struct {
	products map[string]*Product
	mu       sync.RWMutex
}

// NewInMemoryProductRepository creates a product repository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*Product),
	}
}

func (r *InMemoryProductRepository) FindByID(id string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (r *InMemoryProductRepository) FindByCategory(category string) ([]*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*Product, 0)
	for _, product := range r.products {
		if product.Category == category {
			products = append(products, product)
		}
	}
	return products, nil
}

func (r *InMemoryProductRepository) FindInStock() ([]*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*Product, 0)
	for _, product := range r.products {
		if product.InStock {
			products = append(products, product)
		}
	}
	return products, nil
}

func (r *InMemoryProductRepository) Save(product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return errors.New("product already exists")
	}
	r.products[product.ID] = product
	return nil
}

func (r *InMemoryProductRepository) Update(product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}
	r.products[product.ID] = product
	return nil
}

func (r *InMemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(r.products, id)
	return nil
}

// UserService demonstrates using repositories
type UserService struct {
	repo UserRepository
}

// NewUserService creates a user service
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(id, name, email string, age int) error {
	// Check if email exists
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return fmt.Errorf("email %s already registered", email)
	}

	user := &User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}
	return s.repo.Save(user)
}

// GetUser retrieves a user
func (s *UserService) GetUser(id string) (*User, error) {
	return s.repo.FindByID(id)
}

// UpdateUserEmail updates user email
func (s *UserService) UpdateUserEmail(id, newEmail string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	user.Email = newEmail
	return s.repo.Update(user)
}
