package domain

type Restaurant struct {
	ID       int
	Name     string
	OwnerID  int
	ImageURL string
}

func NewRestaurant(id int, name string, ownerID int, imageUrl string) Restaurant {
	return Restaurant{
		ID:       id,
		Name:     name,
		OwnerID:  ownerID,
		ImageURL: imageUrl,
	}
}

func (r *Restaurant) Validate() bool {
	if r.Name == "" || r.OwnerID <= 0 {
		return false
	}
	return true
}
