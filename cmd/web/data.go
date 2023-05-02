// Filename: cmd/web/data.go
package main

import (
	"github.com/abelwhite/poll/internal/models"
)

type templateData struct {
	Question  *models.Question
	Flash     string //flash is the key
	CSRFToken string
}
