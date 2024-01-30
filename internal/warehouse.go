package internal

// WarehouseAttributes is a struct that contains the attributes of a warehouse
type WarehouseAttributes struct {
	// Name is the name of the Warehouse
	Name string
	// Address is the quantity of the Warehouse
	Address string
	// Telephone is the code value of the Warehouse
	Telephone string
	//Capacity is the capacity of the Warehouse
	Capacity int
}

// Warehouse is a struct that contains the attributes of a Warehouse
type Warehouse struct {
	// Id is the unique identifier of the Warehouse
	Id int
	// WarehouseAttributes is the attributes of the Warehouse
	WarehouseAttributes
}

type WarehouseProductsCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
