package store

import "go-cashier-api/model"

// In-memory product data
var categories = []model.Category{
	{ID: 1, Name: "Sayuran", Description: "Sayuran segar dari petani terbaik"},
	{ID: 2, Name: "Buah-buahan", Description: "Buah-buahan pilihan dari lahan petani yang bersih"},
	{ID: 3, Name: "Tanaman Obat", Description: "Tanaman Obat yang memiliki khasiat yang manjur"},
}

// Query functions
func GetAllCategories() []model.Category {
	return categories
}

func GetCategoryByID(id int) (model.Category, bool) {
	// find category by ID
	for _, p := range categories {
		// if ID matches, return category
		if p.ID == id {
			// found, return true
			return p, true
		}
	}
	// if not found, return false
	return model.Category{}, false
}

// Command functions
func CreateCategory(c model.Category) model.Category {
	// assign new ID
	c.ID = len(categories) + 1
	// append to categories slice
	categories = append(categories, c)
	// return created category
	return c
}

func UpdateCategory(id int, updated model.Category) (model.Category, bool) {
	// find category by ID
	for i, c := range categories {
		if c.ID == id {
			updated.ID = id
			categories[i] = updated
			return updated, true
		}
	}
	return model.Category{}, false
}

func DeleteCategory(id int) bool {
	// find product by ID
	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			return true
		}
	}
	return false
}
