// Filename: internal/models/options.go
package models

import (
	"context"
	"database/sql"
	"time"
)

type Options struct {
	OptionID  int64
	Option_1  string
	Option_2  string
	Option_3  string
	Option_4  string
	CreatedAt time.Time
}

type OptionsModel struct {
	DB *sql.DB
}

func (m *OptionsModel) Insert(option_1 string, option_2 string, option_3 string, option_4 string) (int64, error) {
	var id int64
	// build the query
	statement := `
	             INSERT INTO options(option_1, option_2, option_3, option_4)
							 VALUES($1, $2, $3, $4)
							 RETURNING options_id
	             `
	// build a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// write options to the database
	err := m.DB.QueryRowContext(ctx, statement, option_1, option_2, option_3, option_4).Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil
}
