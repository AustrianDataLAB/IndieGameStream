package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"log"
	"net/http"
	"os"
)

type IAuthService interface {
	Authorize(_ *gin.Context)
}

type authService struct {
}

func (_ authService) Authorize(c *gin.Context) {
	//Check if OAUTH_CLIENT has been set
	if len(os.Getenv("OAUTH_CLIENT")) == 0 {
		//if not set & we are in production mode, abort
		if os.Getenv("GIN_MODE") == "release" {
			log.Fatalf("OAUTH_CLIENT is not set, cannot authorize requests")
		}
		//otherwise we are in debug mode and we can accept it
		log.Println("OAUTH_CLIENT is not set, cannot authorize requests")
		c.Set("subject", "")
		return
	}

	tokenString := c.GetHeader("Authorization")
	payload, err := idtoken.Validate(context.Background(), tokenString, os.Getenv("OAUTH_CLIENT"))

	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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
