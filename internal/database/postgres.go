package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// DB wraps a PostgreSQL connection pool.
type DB struct {
	pool *sql.DB
}

// New connects to PostgreSQL and returns a DB instance.
func New(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &DB{pool: db}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() error {
	return db.pool.Close()
}

// CreateUser inserts a new user and returns the created record.
func (db *DB) CreateUser(ctx context.Context, user types.User) (*types.User, error) {
	query := `
		INSERT INTO users (id, username, email, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, role, created_at
	`
	row := db.pool.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, user.Role, user.CreatedAt)
	var created types.User
	if err := row.Scan(&created.ID, &created.Username, &created.Email, &created.Role, &created.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &created, nil
}

// ListUsers returns all users.
func (db *DB) ListUsers(ctx context.Context) ([]types.User, error) {
	query := `SELECT id, username, email, role, created_at FROM users ORDER BY created_at DESC`
	rows, err := db.pool.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()
	var users []types.User
	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// InitSchema creates the required tables if they don't exist.
func (db *DB) InitSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		role TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS incidents (
		id TEXT PRIMARY KEY,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		type TEXT NOT NULL,
		severity TEXT NOT NULL,
		status TEXT NOT NULL,
		summary TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS compliance_monitors (
		monitor_id TEXT PRIMARY KEY,
		framework TEXT NOT NULL,
		status TEXT NOT NULL,
		last_run TIMESTAMP WITH TIME ZONE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS agent_memory (
		key TEXT PRIMARY KEY,
		value JSONB NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL
	);
	`
	if _, err := db.pool.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("failed to init schema: %w", err)
	}
	log.Println("database schema initialized")
	return nil
}
