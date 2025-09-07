package domain

type MenuItem struct {
	ID           int
	Name         string
	Price        float64
	Available    bool
	RestaurantID int
}

func NewMenuItem(id int, name string, price float64, available bool) MenuItem {
	return MenuItem{
		ID:        id,
		Name:      name,
		Price:     price,
		Available: available,
	}
}

func (m *MenuItem) Validate() bool {
	if m.Name == "" || m.Price < 0 || m.RestaurantID <= 0 {
		return false
	}
	return true
}

func (m *MenuItem) IsAvailable() bool {
	return m.Available
}
