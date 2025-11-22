package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCompany(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}
