// internal/repository/postgres.go
package repository

import (
	"context"
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/dagherghinescu/companies/internal/models"
)

// postgresRepo implements CompanyRepository using Postgres + Squirrel
type postgresRepo struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

// NewPostgresRepo creates a new Postgres repository instance
func NewPostgresRepo(db *sql.DB) Company {
	return &postgresRepo{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create inserts a new company record
func (r *postgresRepo) Create(ctx context.Context, c *models.Company) error {
	query := r.sb.Insert("companies").
		Columns("id", "name", "description", "amount_of_employees", "registered", "type").
		Values(c.ID, c.Name, c.Description, c.AmountEmployees, c.Registered, c.Type)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}

// GetByID retrieves a company by ID
func (r *postgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	query := r.sb.Select("id", "name", "description", "amount_of_employees", "registered", "type").
		From("companies").
		Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	var c models.Company
	err = row.Scan(&c.ID, &c.Name, &c.Description, &c.AmountEmployees, &c.Registered, &c.Type)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Patch updates only the specified columns in updates for the company with id.
func (r *postgresRepo) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	q := r.sb.Update("companies")
	for col, val := range updates {
		q = q.Set(col, val)
	}
	q = q.Where(sq.Eq{"id": id})

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a company by ID
func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := r.sb.Delete("companies").
		From("companies").
		Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}

// EnsureAdminUser creates the admin user if it doesn't exist.
func EnsureAdminUser(ctx context.Context, db *sql.DB, username, password string) error {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		log.Println("Admin user already exists")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	id := uuid.New()
	_, err = db.ExecContext(ctx,
		"INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)",
		id, username, string(hash),
	)
	if err != nil {
		return err
	}

	log.Println("Admin user created with username:", username)
	return nil
}
