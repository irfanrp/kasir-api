package repositories

import (
	"database/sql"
	"errors"
	"kasir-api-irfan/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), 
		       COALESCE(c.name, '') as category_name, COALESCE(c.description, '') as category_description
		FROM product p
		LEFT JOIN category c ON p.category_id = c.id
	`

	args := []interface{}{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var catID int
		var catName string
		var catDesc string
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID, &catName, &catDesc)
		if err != nil {
			return nil, err
		}

		p.CategoryID = catID
		if catID != 0 {
			p.Category = &models.Category{
				ID:          catID,
				Name:        catName,
				Description: catDesc,
			}
		}

		products = append(products, p)
	}
	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), 
		       COALESCE(c.name, '') as category_name, COALESCE(c.description, '') as category_description
		FROM product p
		LEFT JOIN category c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var p models.Product
	var catID int
	var catName string
	var catDesc string
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID, &catName, &catDesc)
	if err == sql.ErrNoRows {
		return models.Product{}, errors.New("product not found")
	}
	if err != nil {
		return models.Product{}, err
	}

	p.CategoryID = catID
	if catID != 0 {
		p.Category = &models.Category{
			ID:          catID,
			Name:        catName,
			Description: catDesc,
		}
	}

	return p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("product not found")
	}
	return err
}
