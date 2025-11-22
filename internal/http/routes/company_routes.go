package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/dagherghinescu/companies/internal/app"
	"github.com/dagherghinescu/companies/internal/http/handlers"
	"github.com/dagherghinescu/companies/internal/http/middleware"
)

func RegisterCompanyRoutes(r *gin.Engine, app *app.App, jwtCfg *middleware.JWTConfig, db *sql.DB) {
	auth := r.Group("/", middleware.JWTMiddleware(jwtCfg))
	{
		auth.POST("/companies", handlers.CreateCompany(app))
		auth.PATCH("/companies/:id", handlers.UpdateCompany(app))
		auth.DELETE("/companies/:id", handlers.DeleteCompany(app))
	}

	r.GET("/companies/:id", handlers.GetCompany(app))

	r.POST("/login", handlers.LoginHandler(db, jwtCfg.Secret))
}
