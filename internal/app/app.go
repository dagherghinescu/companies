package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/dagherghinescu/companies/internal/kafka"
	"github.com/dagherghinescu/companies/internal/models"
	"github.com/dagherghinescu/companies/internal/repository"
)

type App struct {
	Logger   *zap.Logger
	DB       repository.Company
	Producer kafka.ProducerInterface
}

// New creates a new App instance
func New(logger *zap.Logger, db repository.Company, producer kafka.ProducerInterface) *App {
	return &App{
		Logger:   logger,
		DB:       db,
		Producer: producer,
	}
}

// CreateCompany creates a new company
func (a *App) CreateCompany(ctx context.Context, c *models.Company) error {
	err := a.DB.Create(ctx, c)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			a.Logger.Debug("create failed - unique violation",
				zap.String("company_name", *c.Name),
				zap.String("detail", pqErr.Detail),
			)

			return ErrCompanyAlreadyExists
		}

		return err
	}

	event := map[string]interface{}{
		"id":     c.ID.String(),
		"name":   *c.Name,
		"action": "created",
	}

	err = a.Producer.Publish(ctx, c.ID.String(), event)
	if err != nil {
		return err
	}

	return nil
}

// GetCompany retrieves a company by ID
func (a *App) GetCompany(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	company, err := a.DB.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		if errors.Is(err, ErrCompanyNotFound) {
			return nil, ErrCompanyNotFound
		}

		return nil, err
	}

	if company == nil {
		return nil, ErrCompanyNotFound
	}

	return company, nil
}

// UpdateCompany updates an existing company
func (a *App) PatchCompany(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	if len(fields) > 0 {
		return nil
	}

	err := a.DB.Patch(ctx, id, fields)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrCompanyNotFound) {
			return ErrCompanyNotFound
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrCompanyAlreadyExists
		}

		return err
	}

	event := map[string]interface{}{
		"id":     id.String(),
		"action": "updated",
		"fields": fields,
	}

	if err := a.Producer.Publish(ctx, id.String(), event); err != nil {
		return err
	}

	return nil
}

// DeleteCompany deletes a company by ID
func (a *App) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	if err := a.DB.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrCompanyNotFound) {
			return ErrCompanyNotFound
		}

		return err
	}

	event := map[string]interface{}{
		"id":     id.String(),
		"action": "deleted",
	}

	if err := a.Producer.Publish(ctx, id.String(), event); err != nil {
		return err
	}

	return nil
}
