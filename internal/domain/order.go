package domain

import "slices"

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

func NewOrder(id int, customerID int, restaurantID int) *Order {
	if customerID <= 0 || restaurantID <= 0 {
		return nil
	}
	return &Order{
		ID:           id,
		CustomerID:   customerID,
		RestaurantID: restaurantID,
		OrderItems:   []OrderItem{},
	}
}

func (o *Order) Validate(menuItems map[int]MenuItem) bool {
	if o.CustomerID <= 0 || o.RestaurantID <= 0 || len(o.OrderItems) == 0 {
		return false
	}
	for _, item := range o.OrderItems {
		if mi, exists := menuItems[item.MenuItemID]; !exists ||
			item.Quantity <= 0 ||
			!mi.IsAvailable() {
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

func (o *Order) RemoveItem(menuItemID int, quantity int) {
	if quantity <= 0 {
		return
	}
	i := slices.IndexFunc(o.OrderItems, func(item OrderItem) bool {
		return item.MenuItemID == menuItemID
	})
	if i == -1 {
		return
	}
	if o.OrderItems[i].Quantity < quantity {
		o.OrderItems = append(o.OrderItems[:i], o.OrderItems[i+1:]...)
		return
	}
	o.OrderItems[i].Quantity -= quantity
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
