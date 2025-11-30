package httpcon

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: wrap with auth user
func getAllUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"response": "Admin Nonser",
	})
}
