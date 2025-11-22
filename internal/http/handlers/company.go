package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/dagherghinescu/companies/internal/app"
	"github.com/dagherghinescu/companies/internal/models"
)

// GetCompany returns a handler that retrieves a company by its UUID.
// It reads the company ID from the request path, validates it,
// calls the application service, and responds with the company data
// or an appropriate HTTP error.
func GetCompany(appl *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		// Validate UUID
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company id"})
			return
		}

		company, err := appl.GetCompany(c.Request.Context(), id)
		if err != nil {
			switch err {
			case app.ErrCompanyNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			default:
				appl.Logger.Error("get failed", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, company)
	}
}

// CreateCompany returns a handler that creates a new company.
// It binds and validates the incoming JSON payload, delegates
// creation to the application service, and responds with the
// created resource or an error if validation or creation fails.
func CreateCompany(appl *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.Company
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input.ID = uuid.New()
		if err := appl.CreateCompany(c.Request.Context(), &input); err != nil {
			if err == app.ErrCompanyAlreadyExists {
				c.JSON(http.StatusConflict, gin.H{"error": "company with that name already exists"})
				return
			}
			appl.Logger.Error("error creating the company", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		c.JSON(http.StatusCreated, input)
	}
}

// UpdateCompany returns a handler for fully updating a company resource.
// It parses the UUID from the path, binds the JSON body, and passes
// the updated data to the application service.
func UpdateCompany(appl *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company id"})
			return
		}

		var input models.Company
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updates := make(map[string]interface{})

		if input.Name != nil {
			updates["name"] = input.Name
		}
		if input.Description != nil {
			updates["description"] = input.Description
		}
		if input.AmountEmployees != nil {
			updates["amount_of_employees"] = input.AmountEmployees
		}
		if input.Registered != nil {
			updates["registered"] = input.Registered
		}
		if input.Type != nil {
			updates["type"] = input.Type
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		err = appl.PatchCompany(c.Request.Context(), id, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			appl.Logger.Error("error patching company", zap.Error(err))
			return
		}

		c.JSON(http.StatusOK, updates)
	}
}

// DeleteCompany returns a handler that deletes a company by ID.
// It expects the company UUID as a path parameter.
func DeleteCompany(appl *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company id"})
			return
		}

		if err := appl.DeleteCompany(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
