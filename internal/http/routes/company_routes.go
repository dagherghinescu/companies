package routes

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/dagherghinescu/companies/internal/app"
	"github.com/dagherghinescu/companies/internal/http/handlers"
	"github.com/dagherghinescu/companies/internal/http/middleware"
	"github.com/dagherghinescu/companies/internal/models"
)

func RegisterCompanyRoutes(r *gin.Engine, app *app.App, jwtCfg *middleware.JWTConfig) {
	auth := r.Group("/", middleware.JWTMiddleware(jwtCfg))
	{
		auth.POST("/companies", handlers.CreateCompany(app))
		auth.PATCH("/companies/:id", handlers.UpdateCompany(app))
		auth.DELETE("/companies/:id", handlers.DeleteCompany(app))
	}

	r.GET("/companies/:id", handlers.GetCompany(app))

	// NOTE: this implementation was done only to make it easier to obtain a token.
	// Ideally this should be handled by a login/register service and users stored
	// in a database.
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	users := map[string]models.User{
		"admin": {ID: "1", Username: "admin", Password: string(hash)},
	}
	r.POST("/login", handlers.LoginHandler(jwtCfg.Secret, users))
}
