package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"log"
	"net/http"
	"strings"
)

type IAuthService interface {
	Authorize(_ *gin.Context)
}

type authService struct {
}

func (_ authService) Authorize(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	payload, err := idtoken.Validate(context.Background(), tokenString, "516825360638-ai7mibm97c1i5o66l18iqlfuqffl1dba.apps.googleusercontent.com")

	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}
	if payload == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("subject", payload.Subject)

}

func AuthService() IAuthService {
	return &authService{}
}
