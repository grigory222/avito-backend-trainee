package repo

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// AddProduct добавление товара в приемку
func (pr *ProductRepository) AddProduct(productType, receptionId string) (models.Product, error) {
	sql := "INSERT INTO products (product_type, reception_id) VALUES ($1, $2) RETURNING id, date_time, product_type, reception_id"
	product := models.Product{}
	err := pr.db.QueryRow(sql, productType, receptionId).Scan(
		&product.Id,
		&product.DateTime,
		&product.Type,
		&product.ReceptionId,
	)
	if err != nil {
		return models.Product{}, models.ErrDBInsert
	}
	return product, nil
}

// GetLastProduct получение последнего товара
func (pr *ProductRepository) GetLastProduct(receptionId string) (models.Product, error) {
	sql := "SELECT * FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1"
	product := models.Product{}
	err := pr.db.QueryRow(sql, receptionId).Scan(
		&product.Id,
		&product.DateTime,
		&product.Type,
		&product.ReceptionId)
	if err != nil {
		return models.Product{}, models.ErrNoProductsInReception
	}
	return product, nil
}

// DeleteProductById удаление товара
func (pr *ProductRepository) DeleteProductById(id string) error {
	sql := "DELETE FROM products WHERE id = $1"
	_, err := pr.db.Exec(sql, id)
	if err != nil {
		return err
	}
	return nil
}

// получение товаров по receptionIds
// GetProductsByReceptionIds(recIds []string, startDate, endDate *time.Time) ([]models.Product, error)
