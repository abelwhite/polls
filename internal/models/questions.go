// Filename: internal/models/questions.go
package models

import (
	"context"
	"database/sql"
	"time"
)

// The Question model will represent a single question in our questions table
type Question struct {
	QuestionID int64
	Body string
	CreatedAt time.Time
}
// The QuestionModel type will encapsulate the
// DB connection pool that will be initialized
// in the main() function
type QuestionModel struct {
	DB *sql.DB
}

// The Insert() function stores a question into the questions table
func (m *QuestionModel) Insert(body string) (int64, error) {
	// id will be used to stored the unique identifier returned by
	// PostgreSQL after adding the row to the table
	var id int64
	statement := 
	            `
							INSERT INTO questions(body)
							VALUES($1)
							RETURNING question_id
	            `
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, statement, body).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *QuestionModel) Get() (*Question, error) {
	var q Question

	statement := 
	            `
							SELECT question_id, body
							FROM questions
							ORDER BY RANDOM()
							LIMIT 1
	            `
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, statement).Scan(&q.QuestionID, &q.Body)
	if err != nil {
		return nil, err
	}
	return &q, nil
}