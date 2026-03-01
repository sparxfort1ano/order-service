package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var ErrOrderNotFound = errors.New("order not found")

type PostgresRepository struct {
	Db *sql.DB
}

func (r *PostgresRepository) SaveOrder(ctx context.Context, o *Order) error {
	jsonData, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := `
        INSERT INTO orders (
            order_uid, track_number, entry, data, 
            locale, internal_signature, customer_id, delivery_service, 
            shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        ON CONFLICT (order_uid) DO NOTHING;
    `

	if _, err = r.Db.ExecContext(
		ctx,
		query,
		o.OrderUid,
		o.TrackNumber,
		o.Entry,
		jsonData, // JSONB
		o.Locale,
		o.InternalSignature,
		o.CustomerID,
		o.DeliveryService,
		o.Shardkey,
		o.SmID,
		o.DateCreated,
		o.OofShard,
	); err != nil {
		return fmt.Errorf("failed to insert into database: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetOrderById(ctx context.Context, orderUID string) (*Order, error) {
	query := `
	SELECT data FROM orders WHERE order_uid = $1`

	row := r.Db.QueryRowContext(ctx, query, orderUID)

	var jsonData []byte
	if err := row.Scan(&jsonData); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to scan data: %w", err)
	}

	var o Order
	if err := json.Unmarshal(jsonData, &o); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &o, nil
}

// Select all orders from database
func (r *PostgresRepository) GetAllOrders(ctx context.Context) ([]*Order, error) {
	queryOrders := `SELECT * FROM orders`
	rows, err := r.Db.QueryContext(ctx, queryOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to select from database: %w", err)
	}
	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var o Order
		var jsonData []byte

		if err := rows.Scan(&o.OrderUid,
			&o.TrackNumber,
			&o.Entry,
			&jsonData,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerID,
			&o.DeliveryService,
			&o.Shardkey,
			&o.SmID,
			&o.DateCreated,
			&o.OofShard,
		); err != nil {
			return nil, fmt.Errorf("failed to read the rows: %w", err)
		}

		if err := json.Unmarshal(jsonData, &o); err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read the rows: %w", err)
	}

	return orders, nil
}

// Migration mechanism
func (r *PostgresRepository) Migrate(path string) error {
	query, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if _, err := r.Db.Exec(string(query)); err != nil {
		return fmt.Errorf("failed to insert into database: %w", err)
	}

	return nil
}

// Close connection with database
func (r *PostgresRepository) Close() error {
	err := r.Db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

// Postgres repository init
func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Database communication check
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{
		Db: db,
	}, nil
}
