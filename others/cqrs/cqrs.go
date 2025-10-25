// Package cqrs demonstrates the CQRS pattern.
// It separates read and write operations using different models and potentially
// different data stores, optimizing each for their specific purpose.
package cqrs

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Command interface for write operations
type Command interface {
	CommandName() string
}

// Query interface for read operations
type Query interface {
	QueryName() string
}

// Commands (Write Side)

type CreateProductCommand struct {
	ProductID string
	Name      string
	Price     float64
	Category  string
}

func (c *CreateProductCommand) CommandName() string { return "CreateProduct" }

type UpdateProductPriceCommand struct {
	ProductID string
	NewPrice  float64
}

func (c *UpdateProductPriceCommand) CommandName() string { return "UpdateProductPrice" }

type DeleteProductCommand struct {
	ProductID string
}

func (c *DeleteProductCommand) CommandName() string { return "DeleteProduct" }

// Queries (Read Side)

type GetProductByIDQuery struct {
	ProductID string
}

func (q *GetProductByIDQuery) QueryName() string { return "GetProductByID" }

type GetProductsByCategoryQuery struct {
	Category string
}

func (q *GetProductsByCategoryQuery) QueryName() string { return "GetProductsByCategory" }

type GetAllProductsQuery struct{}

func (q *GetAllProductsQuery) QueryName() string { return "GetAllProducts" }

// Write Model (Command Side)

type ProductWriteModel struct {
	ID          string
	Name        string
	Price       float64
	Category    string
	Version     int
	LastUpdated time.Time
}

// Read Model (Query Side) - Optimized for reads

type ProductReadModel struct {
	ID           string
	Name         string
	Price        float64
	Category     string
	DisplayPrice string // Pre-formatted
	Searchable   string // Pre-computed search string
}

type ProductListItem struct {
	ID       string
	Name     string
	Price    string
	Category string
}

// Write Store (Normalized)

type WriteStore struct {
	products map[string]*ProductWriteModel
	mu       sync.RWMutex
}

func NewWriteStore() *WriteStore {
	return &WriteStore{
		products: make(map[string]*ProductWriteModel),
	}
}

func (ws *WriteStore) Save(product *ProductWriteModel) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	product.Version++
	product.LastUpdated = time.Now()
	ws.products[product.ID] = product
	return nil
}

func (ws *WriteStore) Get(id string) (*ProductWriteModel, error) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	product, exists := ws.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (ws *WriteStore) Delete(id string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	delete(ws.products, id)
	return nil
}

// Read Store (Denormalized, optimized for queries)

type ReadStore struct {
	products         map[string]*ProductReadModel
	productsByCategory map[string][]*ProductListItem
	mu               sync.RWMutex
}

func NewReadStore() *ReadStore {
	return &ReadStore{
		products:           make(map[string]*ProductReadModel),
		productsByCategory: make(map[string][]*ProductListItem),
	}
}

func (rs *ReadStore) Save(product *ProductReadModel) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.products[product.ID] = product

	// Update category index
	listItem := &ProductListItem{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.DisplayPrice,
		Category: product.Category,
	}

	if _, exists := rs.productsByCategory[product.Category]; !exists {
		rs.productsByCategory[product.Category] = make([]*ProductListItem, 0)
	}
	rs.productsByCategory[product.Category] = append(rs.productsByCategory[product.Category], listItem)
}

func (rs *ReadStore) GetByID(id string) (*ProductReadModel, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	product, exists := rs.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (rs *ReadStore) GetByCategory(category string) []*ProductListItem {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	return rs.productsByCategory[category]
}

func (rs *ReadStore) GetAll() []*ProductListItem {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	all := make([]*ProductListItem, 0)
	for _, products := range rs.productsByCategory {
		all = append(all, products...)
	}
	return all
}

// Command Bus (Write Side)

type CommandHandler interface {
	Handle(cmd Command) error
}

type CommandBus struct {
	handlers map[string]CommandHandler
}

func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

func (cb *CommandBus) Register(cmdName string, handler CommandHandler) {
	cb.handlers[cmdName] = handler
}

func (cb *CommandBus) Send(cmd Command) error {
	handler, exists := cb.handlers[cmd.CommandName()]
	if !exists {
		return fmt.Errorf("no handler for command %s", cmd.CommandName())
	}
	return handler.Handle(cmd)
}

// Query Bus (Read Side)

type QueryHandler interface {
	Handle(query Query) (interface{}, error)
}

type QueryBus struct {
	handlers map[string]QueryHandler
}

func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

func (qb *QueryBus) Register(queryName string, handler QueryHandler) {
	qb.handlers[queryName] = handler
}

func (qb *QueryBus) Execute(query Query) (interface{}, error) {
	handler, exists := qb.handlers[query.QueryName()]
	if !exists {
		return nil, fmt.Errorf("no handler for query %s", query.QueryName())
	}
	return handler.Handle(query)
}

// Command Handlers

type CreateProductHandler struct {
	writeStore *WriteStore
	readStore  *ReadStore
}

func (h *CreateProductHandler) Handle(cmd Command) error {
	createCmd := cmd.(*CreateProductCommand)

	// Save to write store
	writeModel := &ProductWriteModel{
		ID:       createCmd.ProductID,
		Name:     createCmd.Name,
		Price:    createCmd.Price,
		Category: createCmd.Category,
	}
	if err := h.writeStore.Save(writeModel); err != nil {
		return err
	}

	// Project to read store
	readModel := &ProductReadModel{
		ID:           createCmd.ProductID,
		Name:         createCmd.Name,
		Price:        createCmd.Price,
		Category:     createCmd.Category,
		DisplayPrice: fmt.Sprintf("$%.2f", createCmd.Price),
		Searchable:   createCmd.Name + " " + createCmd.Category,
	}
	h.readStore.Save(readModel)

	return nil
}

// Query Handlers

type GetProductByIDHandler struct {
	readStore *ReadStore
}

func (h *GetProductByIDHandler) Handle(query Query) (interface{}, error) {
	q := query.(*GetProductByIDQuery)
	return h.readStore.GetByID(q.ProductID)
}

type GetProductsByCategoryHandler struct {
	readStore *ReadStore
}

func (h *GetProductsByCategoryHandler) Handle(query Query) (interface{}, error) {
	q := query.(*GetProductsByCategoryQuery)
	return h.readStore.GetByCategory(q.Category), nil
}

// ProductService demonstrates CQRS in action
type ProductService struct {
	commandBus *CommandBus
	queryBus   *QueryBus
}

func NewProductService(writeStore *WriteStore, readStore *ReadStore) *ProductService {
	commandBus := NewCommandBus()
	queryBus := NewQueryBus()

	// Register command handlers
	commandBus.Register("CreateProduct", &CreateProductHandler{
		writeStore: writeStore,
		readStore:  readStore,
	})

	// Register query handlers
	queryBus.Register("GetProductByID", &GetProductByIDHandler{readStore: readStore})
	queryBus.Register("GetProductsByCategory", &GetProductsByCategoryHandler{readStore: readStore})

	return &ProductService{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (ps *ProductService) CreateProduct(id, name string, price float64, category string) error {
	cmd := &CreateProductCommand{
		ProductID: id,
		Name:      name,
		Price:     price,
		Category:  category,
	}
	return ps.commandBus.Send(cmd)
}

func (ps *ProductService) GetProduct(id string) (*ProductReadModel, error) {
	query := &GetProductByIDQuery{ProductID: id}
	result, err := ps.queryBus.Execute(query)
	if err != nil {
		return nil, err
	}
	return result.(*ProductReadModel), nil
}

func (ps *ProductService) GetProductsByCategory(category string) []*ProductListItem {
	query := &GetProductsByCategoryQuery{Category: category}
	result, _ := ps.queryBus.Execute(query)
	return result.([]*ProductListItem)
}
