package orderRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"go_learn_project_rest_api/modules/orders"
	"go_learn_project_rest_api/modules/orders/orderPatterns"
	"strings"

	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	FindOneOrder(string) (*orders.Order, error)
	FindOrder(*orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(*orders.Order) (string, error)
	UpdateOrder(*orders.Order) error
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

func (r *orderRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := orderPatterns.InsertOrderBuilder(r.db, req)
	orderId, err := orderPatterns.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}
	return orderId, nil
}

func (r *orderRepository) UpdateOrder(req *orders.Order) error {
	query := `
	UPDATE "orders" SET`

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(values, req.Status)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"status" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(values, req.TransferSlip)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"transfer_slip" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(`
	WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("update order failed: %v", err)
	}
	return nil
}
