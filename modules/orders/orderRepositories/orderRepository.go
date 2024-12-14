package orderRepositories

import "github.com/jmoiron/sqlx"

type IOrderRepository interface{}

type orderRepository struct {
	db *sqlx.DB
}

func OrderRepository(db *sqlx.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}
