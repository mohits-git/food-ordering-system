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

func (o *Order) AddItem(menuItemID int, quantity int) {
	if quantity <= 0 {
		return
	}
	for i, item := range o.OrderItems {
		if item.MenuItemID == menuItemID {
			o.OrderItems[i].Quantity += quantity
			return
		}
	}
	o.OrderItems = append(o.OrderItems, OrderItem{MenuItemID: menuItemID, Quantity: quantity})
}

func (o *Order) ClearItems() {
	o.OrderItems = []OrderItem{}
}

func (o *Order) TotalPrice(menuItems map[int]MenuItem) float64 {
	total := 0.0
	for _, item := range o.OrderItems {
		if menuItem, exists := menuItems[item.MenuItemID]; exists && menuItem.IsAvailable() {
			total += menuItem.Price * float64(item.Quantity)
		}
	}
	return total
}
