package models

import (
	"github.com/google/uuid"
)

// CompanyType defines allowed company types
type CompanyType string

const (
	Corporation        CompanyType = "Corporation"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "SoleProprietorship"
)

// Company represents a company entity
type Company struct {
	ID              uuid.UUID    `json:"id" db:"id"`
	Name            *string      `json:"name" db:"name"`
	Description     *string      `json:"description,omitempty" db:"description"`
	AmountEmployees *int         `json:"amount_of_employees" db:"amount_of_employees"`
	Registered      *bool        `json:"registered" db:"registered"`
	Type            *CompanyType `json:"type" db:"type"`
}
