package domain

type Restaurant struct {
	ID      int
	Name    string
	OwnerID int
}

func NewRestaurant(id int, name string, ownerID int) Restaurant {
	return Restaurant{
		ID:      id,
		Name:    name,
		OwnerID: ownerID,
	}
}

func (r *Restaurant) Validate() bool {
	if r.Name == "" || r.OwnerID <= 0 {
		return false
	}
	return true
}
