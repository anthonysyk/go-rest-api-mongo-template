package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var UnauthorizedResponse = gin.H{"message": "unauthorized"}

func Middleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
			return
		}

		_, claims, err := VerifyToken(bearerToken[1], secret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
			return
		}

		// if id is specified in path params, token must contain user id in claims
		id := ctx.Param("id")
		if id != "" && claims["id"] != id {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
			return
		}

		ctx.Next()
	}
}
