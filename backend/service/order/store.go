package order

import (
	"database/sql"
	"github.com/utsavll0/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO public.orders (\"userId\", total, status, address) values ($1, $2, $3, $4)", order.UserId, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO public.order_items (\"orderId\", \"productId\", quantity, price) values ($1, $2, $3, $4)", orderItem.OrderId, orderItem.ProductId, orderItem.Quantity, orderItem.Price)
	return err
}
