package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/dagherghinescu/companies/internal/models"
)

// CompanyRepository defines the contract for interacting with company data.
// Handlers and services should depend on this interface instead of a concrete implementation.
type Company interface {
	Create(ctx context.Context, c *models.Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error)
	Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
