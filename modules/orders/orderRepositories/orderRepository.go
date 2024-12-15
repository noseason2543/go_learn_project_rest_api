package orderRepositories

import (
	"encoding/json"
	"fmt"
	"go_learn_project_rest_api/modules/orders"
	"go_learn_project_rest_api/modules/orders/orderPatterns"

	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	FindOneOrder(string) (*orders.Order, error)
	FindOrder(*orders.OrderFilter) ([]*orders.Order, int)
}

type orderRepository struct {
	db *sqlx.DB
}

func OrderRepository(db *sqlx.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindOneOrder(orderId string) (*orders.Order, error) {
	query := `
	SELECT
		to_jsonb(t)
	FROM (
		SELECT
			o.id,
			o.user_id,
			o.transfer_slip,
			o.status,
			(
				SELECT
					json_agg(pt)
				FROM (
					SELECT
						spo.id,
						spo.qty,
						spo.product
					FROM products_orders spo
					WHERE spo.order_id = o.id
				) AS pt
			) AS products,
			o.address,
			o.contact,
			(
				SELECT
					SUM(COALESCE((po.product ->> 'price')::FLOAT*(po.qty)::FLOAT, 0))
				FROM products_orders po
				WHERE po.order_id = o.id
			) AS total_paid,
			o.created_at,
			o.updated_at
		FROM orders o
		WHERE o.id = $1
	) AS t;`

	orderData := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}
	raw := make([]byte, 0)
	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("get order failed: %v", err)
	}

	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return orderData, nil
}

func (r *orderRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := orderPatterns.FindOrderBuilder(r.db, req)
	engineer := orderPatterns.FindOrderEngineer(builder)
	return engineer.FindOrder(), engineer.CountOrder()
}
