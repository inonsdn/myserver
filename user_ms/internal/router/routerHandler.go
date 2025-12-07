package router

import (
	"fmt"
	"net/http"
	"time"
	dbcon "userms/internal/dbCon"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func NewRouterHandler(cfg *dbcon.DbConfig) *UserRouterHandler {
	sqlCon, err := dbcon.NewSqlCon(cfg)
	if err != nil {
		fmt.Println("Error when router handler", err)
		return nil
	}
	sqlCon.InitializeSchema()
	return &UserRouterHandler{
		localCon: sqlCon,
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
	userId, err := u.localCon.UserAuthentication(loginReq.Username, loginReq.Password)

	fmt.Println("=========", userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// then get user id from db

	// mock user id
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
	fmt.Println("=========", userId)

	// TODO: implement function to query item from local
	// u.localCon.GetUserById(userId)
	userResult, err := u.localCon.GetUserById(userId)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": gin.H{
			"name": userResult.Name,
			"id":   userResult.Id,
		},
	})
}

func (u *UserRouterHandler) createUser(c *gin.Context) {

	userInfo := dbcon.UserInfo{
		Id: uuid.NewString(),
	}

	c.ShouldBindJSON(&userInfo)
	err := u.localCon.CreateUser(userInfo)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": gin.H{
			"id": userInfo.Id,
		},
	})
}
