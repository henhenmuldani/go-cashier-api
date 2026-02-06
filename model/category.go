package model

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryRequestSwagger struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
