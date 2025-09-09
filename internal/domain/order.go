package domain

type Order struct {
	ID           int
	CustomerID   int
	RestaurantID int
	OrderItems   []OrderItem
}

type OrderItem struct {
	MenuItemID int
	Quantity   int
}

func (oi *OrderItem) Validate() bool {
	return oi.MenuItemID > 0 && oi.Quantity > 0
}

func NewOrder(id int, customerID int, restaurantID int) Order {
	return Order{
		ID:           id,
		CustomerID:   customerID,
		RestaurantID: restaurantID,
		OrderItems:   []OrderItem{},
	}
}

func (o *Order) Validate(menuItems map[int]bool) bool {
	if o.CustomerID <= 0 || o.RestaurantID <= 0 || len(o.OrderItems) == 0 {
		return false
	}
	for _, item := range o.OrderItems {
		if available, exists := menuItems[item.MenuItemID]; !exists ||
			!available ||
			!item.Validate() {
			return false
		}
	}
	return true
}
