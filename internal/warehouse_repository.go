package internal

import "errors"

var (
	// ErrRepositoryProductNotFound is returned when a product is not found.
	ErrRepositoryWarehouseNotFound   = errors.New("repository: Warehouse not found")
	ErrRepositoryWarehouseDuplicated = errors.New("repository: Warehouse duplicated")
)

// RepositoryWarehouse is an interface that contains the methods for a Warehouse repository
type RepositoryWarehouse interface {
	// FindById returns a warehouse by its id
	FindById(id int) (w Warehouse, err error)
	// Save saves a warehouse
	Save(w *Warehouse) (err error)
	ReportProducts(id int) (w []WarehouseProductsCount, err error)
	GetAll() (w []Warehouse, err error)
}
