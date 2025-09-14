package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Репозиторий для работы с Postgres
type PostgresRepo struct{ pool *pgxpool.Pool }

// Конструктор репозитория
func NewPostgresRepo(p *pgxpool.Pool) *PostgresRepo { return &PostgresRepo{pool: p} }

// Создание таблицы orders и индекса, если они еще не существуют
func (r *PostgresRepo) Init(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS orders(
  order_uid  TEXT PRIMARY KEY,        -- уникальный id заказа
  payload    JSONB NOT NULL,          -- сам заказ в формате JSON
  created_at TIMESTAMPTZ NOT NULL DEFAULT now() -- дата вставки
);
CREATE INDEX IF NOT EXISTS idx_orders_payload_gin ON orders USING GIN (payload);
`)
	return err
}

// Добавление или обновление заказа в б/д
func (r *PostgresRepo) Upsert(ctx context.Context, o Order) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
INSERT INTO orders(order_uid, payload)
VALUES ($1, $2)
ON CONFLICT(order_uid) DO UPDATE SET payload = EXCLUDED.payload
`, o.OrderUID, b)
	return err
}

// Получение одного заказа по его id
func (r *PostgresRepo) Get(ctx context.Context, id string) (Order, error) {
	var b []byte
	if err := r.pool.QueryRow(ctx, `SELECT payload FROM orders WHERE order_uid=$1`, id).Scan(&b); err != nil {
		return Order{}, err
	}
	var o Order
	return o, json.Unmarshal(b, &o)
}

// Получение списка последних n заказов
func (r *PostgresRepo) All(ctx context.Context, n int) ([]Order, error) {
	rows, err := r.pool.Query(ctx, `SELECT payload FROM orders ORDER BY created_at DESC LIMIT $1`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Order
	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var o Order
		if err := json.Unmarshal(b, &o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
