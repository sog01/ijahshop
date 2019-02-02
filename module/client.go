package module

import (
	"context"

	"github.com/sog01/ijahshop/module/internal"
)

// Client is entity of data client which originally source from storage
type Client struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// GetClientByID is to used fetch client entity based on given ID
func (mod Module) GetClientByID(ctx context.Context, ID int64) (Client, error) {
	var client Client

	clientIntr, err := mod.internal.GetClientByID(ctx, ID)
	if err != nil {
		return Client{}, err
	}

	client = Client{
		ID:      clientIntr.ID,
		Name:    clientIntr.Name,
		Email:   clientIntr.Email,
		Subject: clientIntr.Subject,
		Message: clientIntr.Message,
	}

	return client, nil
}

// CreateClient to insert client entity into internal storage
func (mod Module) CreateClient(ctx context.Context, data Client) error {
	return mod.internal.CreateClient(ctx, internal.Client{
		ID:      data.ID,
		Name:    data.Name,
		Email:   data.Email,
		Subject: data.Subject,
		Message: data.Message,
	})
}
