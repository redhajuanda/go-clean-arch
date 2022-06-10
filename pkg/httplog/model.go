package httplog

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type LogOutgoingRequest struct {
	tableName  struct{} `pg:"httplog.log_outgoing_requests"`
	ID         int
	TraceID    string
	EventName  string
	Endpoint   string
	Request    interface{}
	StatusCode int
	Response   interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type LogIncomingRequest struct {
	tableName struct{} `pg:"httplog.log_incoming_requests"`
	ID        string
	TraceID   string
	EventName string
	Endpoint  string
	Request   interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LogError struct {
	tableName  struct{} `pg:"httplog.log_errors"`
	ID         string
	TraceID    string
	StatusCode int
	Error      string
	Traces     interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// createSchema creates database schema for LogOutgoingRequest and LogIncomingRequest.
func createSchema(db *pg.DB) error {

	_, err := db.Exec(`
	CREATE SCHEMA IF NOT EXISTS httplog;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS httplog.log_outgoing_requests (
		id bigserial NOT NULL PRIMARY KEY,
		trace_id varchar NOT NULL DEFAULT '',
		event_name varchar NOT NULL DEFAULT '',
		endpoint varchar NOT NULL DEFAULT '',
		request varchar NOT NULL DEFAULT '',
		status_code int NOT NULL,
		response varchar NOT NULL DEFAULT '',
		created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS httplog.log_incoming_requests (
		id bigserial NOT NULL PRIMARY KEY,
		trace_id varchar NOT NULL DEFAULT '',
		event_name varchar NOT NULL DEFAULT '',
		endpoint varchar NOT NULL DEFAULT '',
		request varchar NOT NULL DEFAULT '',
		created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS httplog.log_errors (
		id bigserial NOT NULL PRIMARY KEY,
		trace_id varchar NOT NULL DEFAULT '',
		status_code int NOT NULL,
		error varchar NOT NULL DEFAULT '',
		traces varchar NULL,
		created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return err
	}
	return nil
}
