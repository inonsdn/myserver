package router

import (
	"net/http"
	"time"
	dbcon "userms/internal/dbCon"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserRouterHandler struct {
	localCon *dbcon.SqlCon
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken         string `json:"token"`
	ExpirationTimestamp int64  `json:"expiration"`
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "user pong",
	})
}

func NewRouterHandler() *UserRouterHandler {
	// TODO: make local happen
	sqlCon := dbcon.SqlCon{
		// driver: "",
	}
	return &UserRouterHandler{
		localCon: &sqlCon,
	}
}

func (u *UserRouterHandler) login(c *gin.Context) {
	var loginReq loginRequest

	// bind context request to variable that expected username and password
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check that username and password is matched
	// then get user id from db

	// mock user id
	userId := "no11n23"
	loginTime := time.Now()
	expirationTime := loginTime.Add(time.Duration(tokenTimestamp) * time.Second)

	// generate jwt token
	claims := jwt.MapClaims{
		"userId":      userId,
		"expiredUnix": expirationTime.Unix(),
		"authUnix":    loginTime.Unix(),
		"iss":         "user-ms",
	}

	// generate new token with claims map and method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing token with secret
	signingToken, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken:         signingToken,
		ExpirationTimestamp: tokenTimestamp,
	})
}

func (u *UserRouterHandler) getUserInfo(c *gin.Context) {
	userId := c.Request.Header.Get("userId")

	// TODO: implement function to query item from local
	// u.localCon.GetUserById(userId)

	c.JSON(http.StatusOK, gin.H{
		"message": userId,
	})
}
