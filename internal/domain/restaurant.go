package domain

import "slices"

type Restaurant struct {
	ID      int
	Name    string
	OwnerID int
	Menu    []MenuItem
}

func NewRestaurant(id int, name string, ownerID int) Restaurant {
	return Restaurant{
		ID:      id,
		Name:    name,
		OwnerID: ownerID,
		Menu:    []MenuItem{},
	}
}

func (r *Restaurant) Validate() bool {
	if r.Name == "" || r.OwnerID <= 0 {
		return false
	}
	for _, item := range r.Menu {
		if !item.Validate() {
			return false
		}
	}
	return true
}

func (r *Restaurant) AddMenuItem(item MenuItem) {
	existingIndex := slices.IndexFunc(r.Menu, func(mi MenuItem) bool {
		return mi.ID == item.ID
	})
	if existingIndex != -1 {
		r.Menu[existingIndex] = item
		return
	}
	r.Menu = append(r.Menu, item)
}

func (r *Restaurant) RemoveMenuItem(itemID int) {
	r.Menu = slices.DeleteFunc(r.Menu, func(item MenuItem) bool {
		return item.ID == itemID
	})
}
