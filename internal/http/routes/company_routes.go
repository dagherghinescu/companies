package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dagherghinescu/companies/internal/http/handlers"
)

func RegisterCompanyRoutes(r *gin.Engine) {
	r.GET("/companies/:id", handlers.GetCompany)
}
