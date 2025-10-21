package middlewares

import (
	"Gin_Event_App_Arc/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		fmt.Println("UserID from session", session.Get("userID"))
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Authorization header missing"})
			ctx.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Authorization header format"})
			ctx.Abort()
			return

		}
		claims, err := utils.ValidateJWTRSA(tokenParts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Set("userid", claims.UserId)
		ctx.Set("roles", claims.Roles)
		fmt.Println("Idhr tk to aara hai")
		ctx.Next()

	}
}
