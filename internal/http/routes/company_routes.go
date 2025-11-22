package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/dagherghinescu/companies/internal/app"
	"github.com/dagherghinescu/companies/internal/http/handlers"
)

func RegisterCompanyRoutes(r *gin.Engine, app *app.App) {
	r.GET("/companies/:id", handlers.GetCompany(app))
	r.POST("/companies", handlers.CreateCompany(app))
	r.PATCH("/companies/:id", handlers.UpdateCompany(app))
	r.DELETE("/companies/:id", handlers.DeleteCompany(app))
}
