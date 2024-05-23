package product

import (
	"database/sql"
	"fmt"
	"github.com/utsavll0/ecom/types"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM public.products")
	if err != nil {
		return nil, err
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO public.products (name, description, image, price, quantity) values ($1, $2, $3, $4, $5)", product.Name, product.Description, product.Image, product.Price, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetProductsById(productIds []int) ([]types.Product, error) {
	placeholders := generatePlaceholders(len(productIds))
	query := fmt.Sprintf("SELECT * FROM public.products WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(productIds))
	for i, v := range productIds {
		args[i] = v
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var products []types.Product
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func generatePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	placeholders := make([]string, n)
	for i := 0; i < n; i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Description,
		&product.Image,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt)

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4", product.Name, product.Price, product.Quantity, product.ID)
	if err != nil {
		return err
	}
	return nil
}
