package internal

import (
	"context"
)

// Client is entity of client schema
type Client struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Email   string `db:"email"`
	Subject string `db:"subject"`
	Message string `db:"message"`
}

// GetClientByID to fetch client entity from storage based on given ID
func (intr Internal) GetClientByID(ctx context.Context, ID int64) (Client, error) {
	var (
		client Client
		query  string
	)

	query = `SELECT 
				id, 
				name, 
				email, 
				subject, 
				message from client
			WHERE id = ?`

	db := intr.Storage.DB					
	err := db.SelectContext(ctx, &client, db.Rebind(query), ID)
	return client, err
}

// CreateClient to insert client entity into storage
func (intr Internal) CreateClient(ctx context.Context, data Client) error {
	var query string

	query = `INSERT 
				into client 
			(name, email, subject, message) 
				VALUES
			(?, ?, ?, ?)			
			`

	db := intr.Storage.DB
	_, err := db.ExecContext(ctx, db.Rebind(query), data.Name, data.Email, data.Subject, data.Message)
	return err
}
