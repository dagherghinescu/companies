package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/dagherghinescu/companies/internal/models"
	"github.com/dagherghinescu/companies/internal/repository"
)

func TestPostgresRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPostgresRepo(db)

	id := uuid.New()
	name := "Acme Corp"
	description := "Sample"
	employees := 42
	registered := true
	ctype := models.Corporation

	company := &models.Company{
		ID:              id,
		Name:            &name,
		Description:     &description,
		AmountEmployees: &employees,
		Registered:      &registered,
		Type:            &ctype,
	}

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO companies (id,name,description,amount_of_employees,registered,type) VALUES ($1,$2,$3,$4,$5,$6)`)).
		WithArgs(company.ID,
			company.Name,
			company.Description,
			company.AmountEmployees,
			company.Registered,
			company.Type,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), company)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPostgresRepo(db)

	id := uuid.New()
	name := "Acme Corp"
	description := "Sample"
	employees := 42
	registered := true
	ctype := models.Corporation

	rows := sqlmock.NewRows([]string{"id", "name", "description", "amount_of_employees", "registered", "type"}).
		AddRow(id, name, description, employees, registered, ctype)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description, amount_of_employees, registered, type FROM companies WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, id, got.ID)
	require.Equal(t, name, *got.Name)
	require.Equal(t, description, *got.Description)
	require.Equal(t, employees, *got.AmountEmployees)
	require.Equal(t, registered, *got.Registered)
	require.Equal(t, ctype, *got.Type)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepo_Patch(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPostgresRepo(db)

	id := uuid.New()
	updates := map[string]interface{}{
		"name": "New Name",
	}

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE companies SET name = $1 WHERE id = $2`)).
		WithArgs(updates["name"], id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Patch(context.Background(), id, updates)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewPostgresRepo(db)

	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM companies WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
