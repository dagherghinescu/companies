package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/dagherghinescu/companies/internal/app"
	"github.com/dagherghinescu/companies/internal/http/handlers"
	"github.com/dagherghinescu/companies/internal/models"
)

type mockCompanyRepo struct {
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*models.Company, error)
	CreateFn  func(ctx context.Context, c *models.Company) error
	PatchFn   func(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
	DeleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCompanyRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *mockCompanyRepo) Create(ctx context.Context, c *models.Company) error {
	return m.CreateFn(ctx, c)
}
func (m *mockCompanyRepo) Patch(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	return m.PatchFn(ctx, id, fields)
}
func (m *mockCompanyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

// MockProducer does nothing
type mockProducer struct{}

func (m *mockProducer) Publish(_ context.Context, _ string, _ any) error { return nil }
func (m *mockProducer) Close() error                                     { return nil }

func TestGetCompanyHandler(t *testing.T) {
	id := uuid.New()
	company := &models.Company{ID: id, Name: ptrString("Acme")}

	tests := []struct {
		name         string
		param        string
		mockSetup    func(repo *mockCompanyRepo)
		expectedCode int
	}{
		{
			name:  "success",
			param: id.String(),
			mockSetup: func(m *mockCompanyRepo) {
				m.GetByIDFn = func(_ context.Context, _ uuid.UUID) (*models.Company, error) {
					return company, nil
				}
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid UUID",
			param:        "not-a-uuid",
			mockSetup:    func(_ *mockCompanyRepo) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "not found",
			param: id.String(),
			mockSetup: func(m *mockCompanyRepo) {
				m.GetByIDFn = func(_ context.Context, _ uuid.UUID) (*models.Company, error) {
					return nil, app.ErrCompanyNotFound
				}
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "internal error",
			param: id.String(),
			mockSetup: func(m *mockCompanyRepo) {
				m.GetByIDFn = func(_ context.Context, _ uuid.UUID) (*models.Company, error) {
					return nil, errors.New("db error")
				}
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockRepo := &mockCompanyRepo{}
			mockProducer := &mockProducer{}
			tt.mockSetup(mockRepo)

			logger := zap.NewNop()
			appl := app.New(logger, mockRepo, mockProducer)

			router.GET("/companies/:id", handlers.GetCompany(appl))

			req, _ := http.NewRequest(http.MethodGet, "/companies/"+tt.param, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestCreateCompanyHandler(t *testing.T) {
	company := models.Company{Name: ptrString("Acme")}

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func(repo *mockCompanyRepo)
		expectedCode int
	}{
		{
			name: "success",
			body: company,
			mockSetup: func(m *mockCompanyRepo) {
				m.CreateFn = func(_ context.Context, _ *models.Company) error {
					return nil
				}
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON",
			body:         "invalid-json",
			mockSetup:    func(_ *mockCompanyRepo) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "conflict",
			body: company,
			mockSetup: func(m *mockCompanyRepo) {
				m.CreateFn = func(_ context.Context, _ *models.Company) error {
					return app.ErrCompanyAlreadyExists
				}
			},
			expectedCode: http.StatusConflict,
		},
		{
			name: "internal error",
			body: company,
			mockSetup: func(m *mockCompanyRepo) {
				m.CreateFn = func(_ context.Context, _ *models.Company) error {
					return errors.New("db error")
				}
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			mockRepo := &mockCompanyRepo{}
			mockProducer := &mockProducer{}
			tt.mockSetup(mockRepo)

			logger := zap.NewNop()
			appl := app.New(logger, mockRepo, mockProducer)

			router.POST("/companies", handlers.CreateCompany(appl))

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/companies", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

// helper
func ptrString(s string) *string { return &s }
