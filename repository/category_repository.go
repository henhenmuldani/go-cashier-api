package repository

import (
	"database/sql"
	"go-cashier-api/model"
)

type CategoryRepository interface {
	GetAll() ([]model.Category, error)
	GetByID(id int) (*model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) (int64, error) // Return rows affected
	Delete(id int) (int64, error)                   // Return rows affected
}

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

// Query functions
func (repo *CategoryRepositoryImpl) GetAll() ([]model.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID returns a category by its ID
func (repo *CategoryRepositoryImpl) GetByID(id int) (*model.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id = $1"

	var c model.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Command functions
func (repo *CategoryRepositoryImpl) Create(c *model.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	return repo.db.QueryRow(query, c.Name, c.Description).Scan(&c.ID)
}

func (repo *CategoryRepositoryImpl) Update(category *model.Category) (int64, error) {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (repo *CategoryRepositoryImpl) Delete(id int) (int64, error) {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
