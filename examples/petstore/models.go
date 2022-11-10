package petstore

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Address struct {
	City   string `json:"city"`
	State  string `json:"state"`
	Street string `json:"street"`
	Zip    string `json:"zip"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Customer struct {
	Address  []Address `json:"address"`
	ID       int64     `json:"id"`
	Username string    `json:"username"`
}

type Order struct {
	Complete bool        `json:"complete"`
	ID       int64       `json:"id"`
	PetID    int64       `json:"petId"`
	Quantity int64       `json:"quantity"`
	ShipDate string      `json:"shipDate"`
	Status   OrderStatus `json:"status"` // Order Status
}

type Pet struct {
	Category  Category  `json:"category"`
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	PhotoUrls []string  `json:"photoUrls"`
	Status    PetStatus `json:"status"` // pet status in the store
	Tags      []Tag     `json:"tags"`
}

type Pets []Pet

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	ID         int64  `json:"id"`
	LastName   string `json:"lastName"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Username   string `json:"username"`
	UserStatus int64  `json:"userStatus"` // User Status
}

// OrderStatus Order Status
type OrderStatus string

const (
	Approved  OrderStatus = "approved"
	Delivered OrderStatus = "delivered"
	Placed    OrderStatus = "placed"
)

// PetStatus pet status in the store
type PetStatus string

const (
	Available PetStatus = "available"
	Pending   PetStatus = "pending"
	Sold      PetStatus = "sold"
)
