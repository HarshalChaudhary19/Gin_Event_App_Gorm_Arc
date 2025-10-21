package middlewares

import (
	"Gin_Event_App_Arc/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, check := c.Get("roles")
		roles := role.([]*models.Role)
		if !check {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Getting Roles"})
			c.Abort()
			return
		}
		fmt.Println("Roles hai ye", roles)
		isAdmin := false
		for _, r := range roles {
			fmt.Println("Name and Role", r.Id, r.Name)
			if strings.ToLower(r.Name) == "intern" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}

}
