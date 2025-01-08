package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Contact struct {
	bun.BaseModel `bun:"table:contacts"`

	ID       string          `bun:"id,pk"`
	Data     json.RawMessage `bun:"data,type:jsonb"`
	Location string          `bun:"location,unique"`
}

type PostgresDatabase struct {
	db *bun.DB
}

func NewPostgresDatabase(dsn string) (*PostgresDatabase, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return &PostgresDatabase{db: db}, nil
}

func runMigrations(db *bun.DB) error {
	ctx := context.Background()
	
	// Drop the existing table if it exists
	_, err := db.NewDropTable().
		Model((*Contact)(nil)).
		IfExists().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop table: %v", err)
	}
	
	// Create the table with correct schema
	_, err = db.NewCreateTable().
		Model((*Contact)(nil)).
		Exec(ctx)
	return err
}

func (pg *PostgresDatabase) Create(location string, data map[string]interface{}) (bool, string) {
	ctx := context.Background()

	// Check if location already exists
	exists, err := pg.db.NewSelect().
		Model((*Contact)(nil)).
		Where("location = ?", location).
		Exists(ctx)
	if err != nil {
		return false, fmt.Sprintf("Error checking location: %v", err)
	}
	if exists {
		return false, "Location already exists"
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Sprintf("Error marshaling data: %v", err)
	}

	contact := &Contact{
		ID:       location,
		Location: location,
		Data:     jsonData,
	}

	_, err = pg.db.NewInsert().
		Model(contact).
		Exec(ctx)
	if err != nil {
		return false, fmt.Sprintf("Error creating record: %v", err)
	}

	return true, ""
}

func (pg *PostgresDatabase) Read(location string) (bool, string, map[string]interface{}) {
	ctx := context.Background()
	contact := new(Contact)

	err := pg.db.NewSelect().
		Model(contact).
		Where("location = ?", location).
		Scan(ctx)

	if err != nil {
		return false, "Location does not exist", nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(contact.Data, &data); err != nil {
		return false, fmt.Sprintf("Error unmarshaling data: %v", err), nil
	}

	return true, "", data
}

func (pg *PostgresDatabase) Update(location string, data map[string]interface{}) (bool, string) {
	ctx := context.Background()

	// Check if location exists
	exists, err := pg.db.NewSelect().
		Model((*Contact)(nil)).
		Where("location = ?", location).
		Exists(ctx)
	if err != nil {
		return false, fmt.Sprintf("Error checking location: %v", err)
	}
	if !exists {
		return false, "Location does not exist"
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Sprintf("Error marshaling data: %v", err)
	}

	// Create contact with updated data
	contact := &Contact{
		ID:       location,
		Location: location,
		Data:     jsonData,
	}

	// Update using the model
	_, err = pg.db.NewUpdate().
		Model(contact).
		Column("data").
		Where("location = ?", location).
		Exec(ctx)
	if err != nil {
		return false, fmt.Sprintf("Error updating record: %v", err)
	}

	return true, ""
}

func (pg *PostgresDatabase) Delete(location string) (bool, string) {
	ctx := context.Background()

	result, err := pg.db.NewDelete().
		Model((*Contact)(nil)).
		Where("location = ?", location).
		Exec(ctx)
	if err != nil {
		return false, fmt.Sprintf("Error deleting record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Sprintf("Error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return false, "Location does not exist"
	}

	return true, ""
}

func (pg *PostgresDatabase) Close() error {
	return pg.db.Close()
}
